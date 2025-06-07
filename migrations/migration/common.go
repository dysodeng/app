package migration

import (
	"github.com/dysodeng/app/internal/infrastructure/persistence/model/common"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var commonMigration = []*gormigrate.Migration{
	{
		ID: "common_202410231400",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(
				&common.MailConfig{},
				&common.SmsConfig{},
				&common.SmsTemplate{},
			)
			if err != nil {
				return err
			}
			tx.Exec("ALTER TABLE " + (common.SmsConfig{}).TableName() + " COMMENT=\"短信提供商配置表\"")
			tx.Exec("ALTER TABLE " + (common.SmsTemplate{}).TableName() + " COMMENT=\"短信模板表\"")
			tx.Exec("ALTER TABLE " + (common.MailConfig{}).TableName() + " COMMENT=\"邮件服务提供商配置表\"")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&common.MailConfig{},
				&common.SmsConfig{},
				&common.SmsTemplate{},
			)
		},
	},
}
