package model

import (
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	ID   uint   `json:"id" gorm:"primaryKey"` //权限唯一主键
	Name string `json:"name"`                 //权限名称
	// Determine permissions by bit
	//   0:  can see hidden files
	//   1:  can access without password
	//   2:  can add offline download tasks
	//   3:  can mkdir and upload
	//   4:  can rename
	//   5:  can move
	//   6:  can copy
	//   7:  can remove
	//   8:  webdav read
	//   9:  webdav write
	//   10: ftp/sftp login and read
	//   11: ftp/sftp write
	//   12: can read archives
	//   13: can decompress archives
	//   14: dir access control
	Permission  int32          `json:"permission"`
	PathPattern string         `json:"path_pattern"`                              // 目录路径模式
	AllowOp     datatypes.JSON `gorm:"type:json;column:allow_op" json:"allow_op"` //允许的操作upload/download/delete等
	AllowOpInfo AllowOpSlice   `gorm:"-" json:"allow_op_info"`
	CreateTime  time.Time      `json:"create_time"` //创建时间
	UpdateTime  time.Time      `json:"update_time"` //修改时间
}

type AllowOpSlice []string

func (p *Permission) BeforeCreate(db *gorm.DB) (err error) {
	if p.AllowOpInfo != nil {
		p.AllowOp, err = json.Marshal(p.AllowOpInfo)
		if err != nil {
			return
		}
	}
	return nil
}

func (p *Permission) BeforeUpdate(db *gorm.DB) (err error) {
	if p.AllowOpInfo != nil {
		p.AllowOp, err = json.Marshal(p.AllowOpInfo)
		if err != nil {
			return
		}
	}
	return nil
}

func (p *Permission) AfterFind(db *gorm.DB) (err error) {
	p.AllowOpInfo = AllowOpSlice{}
	if len(p.AllowOp) > 0 {
		err = json.Unmarshal(p.AllowOp, &p.AllowOpInfo)
		if err != nil {
			return
		}
	}
	return nil
}

func (p *Permission) CanSeeHides() bool {
	return p.Permission&1 == 1
}

func (p *Permission) CanAccessWithoutPassword() bool {
	return (p.Permission>>1)&1 == 1
}

func (p *Permission) CanAddOfflineDownloadTasks() bool {
	return (p.Permission>>2)&1 == 1
}

func (p *Permission) CanWrite() bool {
	return (p.Permission>>3)&1 == 1
}

func (p *Permission) CanRename() bool {
	return (p.Permission>>4)&1 == 1
}

func (p *Permission) CanMove() bool {
	return (p.Permission>>5)&1 == 1
}

func (p *Permission) CanCopy() bool {
	return (p.Permission>>6)&1 == 1
}

func (p *Permission) CanRemove() bool {
	return (p.Permission>>7)&1 == 1
}

func (p *Permission) CanWebdavRead() bool {
	return (p.Permission>>8)&1 == 1
}

func (p *Permission) CanWebdavManage() bool {
	return (p.Permission>>9)&1 == 1
}

func (p *Permission) CanFTPAccess() bool {
	return (p.Permission>>10)&1 == 1
}

func (p *Permission) CanFTPManage() bool {
	return (p.Permission>>11)&1 == 1
}

func (p *Permission) CanReadArchives() bool {
	return (p.Permission>>12)&1 == 1
}

func (p *Permission) CanDecompress() bool {
	return (p.Permission>>13)&1 == 1
}

func (p *Permission) CanAccessDir() bool {
	return (p.Permission>>14)&1 == 1
}
