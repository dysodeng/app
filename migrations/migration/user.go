package migration

import (
	"github.com/dysodeng/app/internal/infrastructure/persistence/model/user"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var userMigration = []*gormigrate.Migration{
	{
		ID: "user_info_202411091600",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(
				&user.User{},
			)
			if err != nil {
				return err
			}
			tx.Exec("ALTER TABLE " + (user.User{}).TableName() + " COMMENT=\"用户表\"")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&user.User{})
		},
	},
}
