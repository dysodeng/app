package file

import (
	"context"

	fileDomainModel "github.com/dysodeng/app/internal/domain/file/model"
	domainRepository "github.com/dysodeng/app/internal/domain/file/repository"
	fileDataModel "github.com/dysodeng/app/internal/infrastructure/persistence/model/file"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type uploaderRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewUploaderRepository(txManager transactions.TransactionManager) domainRepository.UploaderRepository {
	return &uploaderRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.file.UploaderRepository",
		txManager:         txManager,
	}
}

func (repo *uploaderRepository) CreateMultipartUpload(ctx context.Context, mu *fileDomainModel.MultipartUpload) error {
	dataModel := fileDataModel.MultipartUpload{
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

func (repo *uploaderRepository) FindMultipartUploadByUploadId(ctx context.Context, uploadId string) (*fileDomainModel.MultipartUpload, error) {
	var mu fileDataModel.MultipartUpload
	err := repo.txManager.GetTx(ctx).Where("upload_id=?", uploadId).First(&mu).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return fileDomainModel.MultipartUploadFormModel(&mu), nil
}

// MultipartUploadStatus 分片上传状态设置
func (repo *uploaderRepository) MultipartUploadStatus(ctx context.Context, uploadId string, status uint8) error {
	if err := repo.txManager.GetTx(ctx).Model(&fileDataModel.MultipartUpload{}).
		Where("upload_id=?", uploadId).
		Updates(map[string]interface{}{
			"status": status,
		}).Error; err != nil {
		return err
	}
	return nil
}
