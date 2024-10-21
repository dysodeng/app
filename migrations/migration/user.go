package migration

import (
	"github.com/dysodeng/app/internal/model/common"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var userMigration = []*gormigrate.Migration{
	{
		ID: "2023-02-08T00:00:00",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&common.MailConfig{},
			)
		},
	},
}
