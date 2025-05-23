package data

import (
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func initPermissions() {
	_, err := op.GetPermissionByName("guest")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pg := &model.Permission{
				Name:        "guest",
				Permission:  0x4000, // 14 bit（can access dir）
				PathPattern: "*",
				AllowOpInfo: []string{"upload", "download", "delete"},
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			}
			if err := op.CreatePermission(pg); err != nil {
				panic(err)
			} else {
				utils.Log.Infof("Successfully created the guest permission ")
			}
		}
	}

	_, err = op.GetPermissionByName("admin")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pa := &model.Permission{
				Name:        "admin",
				Permission:  0x70FF, // 0、1、2、3、4、5、6、7、12、13、14 bit
				PathPattern: "",
				AllowOpInfo: []string{"upload", "download", "delete"},
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			}
			if err := op.CreatePermission(pa); err != nil {
				panic(err)
			} else {
				utils.Log.Infof("Successfully created the admin permission ")
			}
		}
	}
}
