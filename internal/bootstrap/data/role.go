package data

import (
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func initRoles() {
	_, err := op.GetRoleByName("guest")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			roleGuest := &model.Role{
				Name:           "guest",
				PermissionInfo: []uint{1},
				CreateTime:     time.Time{},
				UpdateTime:     time.Time{},
			}
			if err := op.CreateRole(roleGuest); err != nil {
				panic(err)
			} else {
				utils.Log.Infof("Successfully created the guest role ")
			}
		}
	}

	_, err = op.GetRoleByName("admin")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			roleAdmin := &model.Role{
				Name:           "admin",
				PermissionInfo: []uint{2},
				CreateTime:     time.Time{},
				UpdateTime:     time.Time{},
			}
			if err := op.CreateRole(roleAdmin); err != nil {
				panic(err)
			} else {
				utils.Log.Infof("Successfully created the admin role ")
			}
		}
	}

}
