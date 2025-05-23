package db

import (
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/pkg/errors"
)

func CreatePermission(p *model.Permission) error {
	return errors.WithStack(db.Create(p).Error)
}

func UpdatePermission(p *model.Permission) error {
	return errors.WithStack(db.Save(p).Error)
}

func GetPermissions(pageIndex, pageSize int) (permissions []model.Permission, count int64, err error) {
	permissionDB := db.Model(&model.Permission{})
	if err := permissionDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get permissions count")
	}
	if err := permissionDB.Order(columnName("id")).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&permissions).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get find permissions")
	}
	return permissions, count, nil
}

func DeletePermissionById(id uint) error {
	return errors.WithStack(db.Delete(&model.Permission{}, id).Error)
}

func GetPermissionById(id uint) (*model.Permission, error) {
	var p model.Permission
	if err := db.First(&p, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get permission by id %d", id)
	}
	return &p, nil
}

func GetPermissionByIds(ids []uint) ([]model.Permission, error) {
	var p []model.Permission
	if err := db.Where("id in (?)", ids).First(&p).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get permission by ids %v", ids)
	}
	return p, nil
}

func GetPermissionByName(name string) (*model.Permission, error) {
	var p model.Permission
	if err := db.Where("name = ?", name).First(&p).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get permission by name %v", name)
	}
	return &p, nil
}
