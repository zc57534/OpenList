package yandex_disk

import (
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
)

type Addition struct {
	RefreshToken   string `json:"refresh_token" required:"true"`
	OrderBy        string `json:"order_by" type:"select" options:"name,path,created,modified,size" default:"name"`
	OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" default:"asc"`
	driver.RootPath
	UseOnlineAPI bool   `json:"use_online_api" default:"true"`
	APIAddress   string `json:"api_url_address" default:"https://api.oplist.org/yandexui/renewapi"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var config = driver.Config{
	Name:        "YandexDisk",
	DefaultRoot: "/",
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &YandexDisk{}
	})
}
