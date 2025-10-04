package file

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/dysodeng/app/internal/domain/file/model"
	"github.com/dysodeng/app/internal/domain/file/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/entity/file"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
)

type uploaderRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewUploaderRepository(txManager transactions.TransactionManager) repository.UploaderRepository {
	return &uploaderRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.file.UploaderRepository",
		txManager:         txManager,
	}
}

func (repo *uploaderRepository) CreateMultipartUpload(ctx context.Context, mu *model.MultipartUpload) error {
	dataModel := file.MultipartUpload{
		UploadID: mu.UploadID,
		FileName: mu.FileName,
		Path:     mu.Path,
		Size:     mu.Size,
		MimeType: mu.MimeType,
		Ext:      mu.Ext,
		Status:   mu.Status,
	}

	tx := repo.txManager.GetTx(ctx)
	if err := tx.Debug().Create(&dataModel).Error; err != nil {
		return err
	}
	mu.ID = dataModel.ID

	return nil
}

func (repo *uploaderRepository) FindMultipartUploadByUploadId(ctx context.Context, uploadId string) (*model.MultipartUpload, error) {
	var mu file.MultipartUpload
	err := repo.txManager.GetTx(ctx).Where("upload_id=?", uploadId).First(&mu).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return repo.multipartUploadFormModel(&mu), nil
}

// MultipartUploadStatus 分片上传状态设置
func (repo *uploaderRepository) MultipartUploadStatus(ctx context.Context, uploadId string, status uint8) error {
	if err := repo.txManager.GetTx(ctx).Model(&file.MultipartUpload{}).
		Where("upload_id=?", uploadId).
		Updates(map[string]interface{}{
			"status": status,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *uploaderRepository) multipartUploadFormModel(mu *file.MultipartUpload) *model.MultipartUpload {
	return &model.MultipartUpload{
		ID:        mu.ID,
		FileName:  mu.FileName,
		Path:      mu.Path,
		Size:      mu.Size,
		MimeType:  mu.MimeType,
		Ext:       mu.Ext,
		UploadID:  mu.UploadID,
		Status:    mu.Status,
		CreatedAt: mu.CreatedAt.Time,
	}
}
