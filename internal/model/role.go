package model

import (
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID             uint              `json:"id" gorm:"primaryKey"`                            //角色唯一主键
	Name           string            `json:"name"`                                            //角色名称
	Permissions    datatypes.JSON    `gorm:"type:json;column:permissions" json:"permissions"` //权限id
	PermissionInfo PermissionIdSlice `gorm:"-" json:"permission_info"`
	CreateTime     time.Time         `json:"create_time"` //创建时间
	UpdateTime     time.Time         `json:"update_time"` //修改时间
}
type PermissionIdSlice []uint

func (r *Role) BeforeCreate(db *gorm.DB) (err error) {
	if r.PermissionInfo != nil {
		r.Permissions, err = json.Marshal(r.PermissionInfo)
		if err != nil {
			return
		}
	}
	return nil
}

func (r *Role) BeforeUpdate(db *gorm.DB) (err error) {
	if r.PermissionInfo != nil {
		r.Permissions, err = json.Marshal(r.PermissionInfo)
		if err != nil {
			return
		}
	}
	return nil
}

func (r *Role) AfterFind(db *gorm.DB) (err error) {
	r.PermissionInfo = PermissionIdSlice{}
	if len(r.Permissions) > 0 {
		err = json.Unmarshal(r.Permissions, &r.PermissionInfo)
		if err != nil {
			return
		}
	}
	return nil
}
