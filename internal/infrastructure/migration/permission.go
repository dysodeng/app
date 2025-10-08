package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/dysodeng/app/internal/infrastructure/persistence/entity/permission"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/model"
)

var permissionMigrations = []*gormigrate.Migration{
	{
		ID: "permission_202510082300",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&permission.Admin{}, &permission.Permission{}, &permission.AdminHasPermission{})
			if err != nil {
				return err
			}
			model.TableComment(tx, db.Driver(), (permission.Admin{}).TableName(), "管理员表")
			model.TableComment(tx, db.Driver(), (permission.Permission{}).TableName(), "管理权限节点表")
			model.TableComment(tx, db.Driver(), (permission.AdminHasPermission{}).TableName(), "管理员权限关联表")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&permission.Admin{}, &permission.Permission{}, &permission.AdminHasPermission{})
		},
	},
}
