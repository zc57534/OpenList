package db

import (
	"errors"
	"github.com/OpenListTeam/OpenList/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FindShareLink 查找分享链接
func FindShareLink(uuid string) (*model.Share, error) {
	var share model.Share
	result := db.First(&share, "uuid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("share link not found")
		}
		return nil, result.Error
	}
	return &share, nil
}

// PushShareLink 插入分享链接
func PushShareLink(share *model.Share) error {
	// 使用 Create 方法插入新记录
	return db.Clauses(clause.OnConflict{
		DoNothing: true, // 如果 UUID 冲突则不做任何操作
	}).Create(share).Error
}

// KillShareLink 删除分享链接
func KillShareLink(uuid string) error {
	// 根据 UUID 删除记录
	return db.Where("uuid = ?", uuid).Delete(&model.Share{}).Error
}

// BatchKillShareLinks 批量删除分享链接（可选扩展）
func BatchKillShareLinks(uuids []string) error {
	return db.Where("uuid IN ?", uuids).Delete(&model.Share{}).Error
}

// UpdateShareLink 更新分享链接（可选扩展）
func UpdateShareLink(uuid string, updates model.Share) error {
	return db.Model(&model.Share{}).Where("uuid = ?", uuid).Updates(updates).Error
}

// FindShareLinks 批量查找分享链接（可选扩展）
func FindShareLinks(uuids []string) ([]model.Share, error) {
	var shares []model.Share
	if err := db.Where("uuid IN ?", uuids).Find(&shares).Error; err != nil {
		return nil, err
	}
	return shares, nil
}
