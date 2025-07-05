package handles

import (
	"bytes"
	"fmt"
	"io"
	stdpath "path"
	"strconv"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/fs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/setting"
	"github.com/OpenListTeam/OpenList/v4/internal/sign"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
)

func Down(c *gin.Context) {
	rawPath := c.MustGet("path").(string)
	filename := stdpath.Base(rawPath)
	storage, err := fs.GetStorage(rawPath, &fs.GetStoragesArgs{})
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if common.ShouldProxy(storage, filename) {
		Proxy(c)
		return
	} else {
		link, _, err := fs.Link(c, rawPath, model.LinkArgs{
			IP:       c.ClientIP(),
			Header:   c.Request.Header,
			Type:     c.Query("type"),
			Redirect: true,
		})
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
		down(c, link)
	}
}

func Proxy(c *gin.Context) {
	rawPath := c.MustGet("path").(string)
	filename := stdpath.Base(rawPath)
	storage, err := fs.GetStorage(rawPath, &fs.GetStoragesArgs{})
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if canProxy(storage, filename) {
		downProxyUrl := storage.GetStorage().DownProxyUrl
		if downProxyUrl != "" {
			_, ok := c.GetQuery("d")
			if !ok {
				URL := fmt.Sprintf("%s%s?sign=%s",
					strings.Split(downProxyUrl, "\n")[0],
					utils.EncodePath(rawPath, true),
					sign.Sign(rawPath))
				c.Redirect(302, URL)
				return
			}
		}
		link, file, err := fs.Link(c, rawPath, model.LinkArgs{
			Header: c.Request.Header,
			Type:   c.Query("type"),
		})
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
		localProxy(c, link, file, storage.GetStorage().ProxyRange)
	} else {
		common.ErrorStrResp(c, "proxy not allowed", 403)
		return
	}
}

func down(c *gin.Context, link *model.Link) {
	if clr, ok := link.MFile.(io.Closer); ok {
		defer func(clr io.Closer) {
			err := clr.Close()
			if err != nil {
				log.Errorf("close link data error: %v", err)
			}
		}(clr)
	}
	var err error
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
	if setting.GetBool(conf.ForwardDirectLinkParams) {
		query := c.Request.URL.Query()
		for _, v := range conf.SlicesMap[conf.IgnoreDirectLinkParams] {
			query.Del(v)
		}
		link.URL, err = utils.InjectQuery(link.URL, query)
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
	}
	c.Redirect(302, link.URL)
}

func localProxy(c *gin.Context, link *model.Link, file model.Obj, proxyRange bool) {
	var err error
	if link.URL != "" && setting.GetBool(conf.ForwardDirectLinkParams) {
		query := c.Request.URL.Query()
		for _, v := range conf.SlicesMap[conf.IgnoreDirectLinkParams] {
			query.Del(v)
		}
		link.URL, err = utils.InjectQuery(link.URL, query)
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
	}
	if proxyRange {
		common.ProxyRange(c, link, file.GetSize())
	}
	Writer := &common.WrittenResponseWriter{ResponseWriter: c.Writer}

	//优先处理md文件
	if utils.Ext(file.GetName()) == "md" && setting.GetBool(conf.FilterReadMeScripts) {
		buf := bytes.NewBuffer(make([]byte, 0, file.GetSize()))
		w := &common.InterceptResponseWriter{ResponseWriter: Writer, Writer: buf}
		err = common.Proxy(w, c.Request, link, file)
		if err == nil && buf.Len() > 0 {
			if c.Writer.Status() < 200 || c.Writer.Status() > 300 {
				c.Writer.Write(buf.Bytes())
				return
			}

			var html bytes.Buffer
			if err = goldmark.Convert(buf.Bytes(), &html); err != nil {
				err = fmt.Errorf("markdown conversion failed: %w", err)
			} else {
				buf.Reset()
				err = bluemonday.UGCPolicy().SanitizeReaderToWriter(&html, buf)
				if err == nil {
					Writer.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))
					Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
					_, err = utils.CopyWithBuffer(Writer, buf)
				}
			}
		}
	} else {
		err = common.Proxy(Writer, c.Request, link, file)
	}
	if err == nil {
		return
	}
	if Writer.IsWritten() {
		log.Errorf("%s %s local proxy error: %+v", c.Request.Method, c.Request.URL.Path, err)
	} else {
		common.ErrorResp(c, err, 500, true)
	}
}

// TODO need optimize
// when can be proxy?
// 1. text file
// 2. config.MustProxy()
// 3. storage.WebProxy
// 4. proxy_types
// solution: text_file + shouldProxy()
func canProxy(storage driver.Driver, filename string) bool {
	if storage.Config().MustProxy() || storage.GetStorage().WebProxy || storage.GetStorage().WebdavProxy() {
		return true
	}
	if utils.SliceContains(conf.SlicesMap[conf.ProxyTypes], utils.Ext(filename)) {
		return true
	}
	if utils.SliceContains(conf.SlicesMap[conf.TextTypes], utils.Ext(filename)) {
		return true
	}
	return false
}
