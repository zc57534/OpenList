package model

type Share struct {
	UUID uint   `json:"share_uuid" gorm:"primaryKey"`
	Date uint   `json:"expire_day"`
	User string `json:"share_user" binding:"required"`
	Path string `json:"share_path" binding:"required"`
	Pass string `json:"share_pass"`
}
