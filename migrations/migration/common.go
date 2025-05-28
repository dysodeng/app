package migration

import (
	common2 "github.com/dysodeng/app/internal/infrastructure/persistence/model/common"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var commonMigration = []*gormigrate.Migration{
	{
		ID: "common_202410231400",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(
				&common2.MailConfig{},
				&common2.SmsConfig{},
				&common2.SmsTemplate{},
			)
			if err != nil {
				return err
			}
			tx.Exec("ALTER TABLE " + (common2.SmsConfig{}).TableName() + " COMMENT=\"短信提供商配置表\"")
			tx.Exec("ALTER TABLE " + (common2.SmsTemplate{}).TableName() + " COMMENT=\"短信模板表\"")
			tx.Exec("ALTER TABLE " + (common2.MailConfig{}).TableName() + " COMMENT=\"邮件服务提供商配置表\"")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&common2.MailConfig{},
				&common2.SmsConfig{},
				&common2.SmsTemplate{},
			)
		},
	},
}
