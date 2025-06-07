package migration

import (
	"github.com/dysodeng/app/internal/infrastructure/persistence/model/file"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var fileMigration = []*gormigrate.Migration{
	{
		ID: "file_202506071400",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&file.File{}, &file.MultipartUpload{}); err != nil {
				return err
			}
			tx.Exec("ALTER TABLE " + (file.File{}).TableName() + " COMMENT=\"文件记录表\"")
			tx.Exec("ALTER TABLE " + (file.MultipartUpload{}).TableName() + " COMMENT=\"文件分片上传记录表\"")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&file.File{}, &file.MultipartUpload{})
		},
	},
}
