package op

import (
	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/internal/model"
)

func CreatePermission(p *model.Permission) error {
	return db.CreatePermission(p)
}

func GetPermissions(pageIndex, pageSize int) (permissions []model.Permission, count int64, err error) {
	return db.GetPermissions(pageIndex, pageSize)
}

func GetPermissionById(id uint) (*model.Permission, error) {
	return db.GetPermissionById(id)
}

func UpdatePermission(p *model.Permission) error {
	return db.UpdatePermission(p)
}

func DeletePermissionById(id uint) error {
	return db.DeletePermissionById(id)
}

func GetPermissionByName(name string) (*model.Permission, error) {
	return db.GetPermissionByName(name)
}
