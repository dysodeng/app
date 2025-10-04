package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/dysodeng/app/internal/infrastructure/persistence/entity/file"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/model"
)

var fileMigrations = []*gormigrate.Migration{
	{
		ID: "file_202510041830",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&file.File{}, &file.MultipartUpload{}); err != nil {
				return err
			}
			model.TableComment(tx, db.Driver(), (file.File{}).TableName(), "文件记录表")
			model.TableComment(tx, db.Driver(), (file.MultipartUpload{}).TableName(), "文件分片上传记录表")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&file.File{}, &file.MultipartUpload{})
		},
	},
}
