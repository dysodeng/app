package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/dysodeng/app/internal/infrastructure/persistence/entity/user"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/model"
)

var userMigrations = []*gormigrate.Migration{
	{
		ID: "user_202510061500",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&user.User{})
			if err != nil {
				return err
			}
			model.TableComment(tx, db.Driver(), (user.User{}).TableName(), "用户表")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&user.User{})
		},
	},
}
