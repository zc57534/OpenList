package op

import (
	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/internal/model"
)

func CreateRole(r *model.Role) error {
	return db.CreateRole(r)
}

func GetRoles(pageIndex, pageSize int) (roles []model.Role, count int64, err error) {
	return db.GetRoles(pageIndex, pageSize)
}

func GetRoleById(id uint) (*model.Role, error) {
	return db.GetRoleById(id)
}

func GetRoleByIds(ids []uint) ([]model.Role, error) {
	return db.GetRoleByIds(ids)
}

func UpdateRole(r *model.Role) error {
	return db.UpdateRole(r)
}

func DeleteRoleById(id uint) error {
	return db.DeleteRoleById(id)
}

func GetPermissionByRoleIds(ids []uint) ([]model.Permission, error) {
	roles, err := db.GetRoleByIds(ids)
	if err != nil {
		return nil, err
	}
	perIds := make([]uint, 0)
	for _, v := range roles {
		perIds = append(perIds, v.PermissionInfo...)
	}
	permissions, err := db.GetPermissionByIds(perIds)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func GetRoleByName(name string) (*model.Role, error) {
	return db.GetRoleByName(name)
}
