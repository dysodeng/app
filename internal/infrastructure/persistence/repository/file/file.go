package file

import (
	"context"
	"strings"

	"github.com/dysodeng/app/internal/domain/file/model"
	fileDomainModel "github.com/dysodeng/app/internal/domain/file/model"
	domainRepository "github.com/dysodeng/app/internal/domain/file/repository"
	fileDataModel "github.com/dysodeng/app/internal/infrastructure/persistence/model/file"
	"github.com/dysodeng/app/internal/infrastructure/persistence/repository"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type fileRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewFileRepository(txManager transactions.TransactionManager) domainRepository.FileRepository {
	return &fileRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.file.FileRepository",
		txManager:         txManager,
	}
}

func (repo *fileRepository) FindList(ctx context.Context, query domainRepository.FileQuery) ([]model.File, int64, error) {
	tx := repo.txManager.GetTx(ctx)

	// 构建查询条件
	db := tx.Debug().Model(&fileDataModel.File{})

	if query.MediaType > 0 {
		db = db.Where("media_type=?", query.MediaType)
	}

	if query.Keyword != "" {
		db = repository.WhereLike(db, "name", query.Keyword)
	}

	if query.StartTime != nil {
		db = db.Where("created_at >= ?", query.StartTime)
	}

	if query.EndTime != nil {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if query.OrderBy != "" {
		orderType := "asc"
		if strings.ToLower(query.OrderType) == "desc" {
			orderType = "desc"
		}
		db = db.Order(query.OrderBy + " " + orderType)
	}

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 查询数据
	var files []fileDataModel.File
	if err := db.Find(&files).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	return fileDomainModel.FileListFromModel(ctx, files), total, nil
}

func (repo *fileRepository) FindByID(ctx context.Context, id uint64) (*model.File, error) {
	tx := repo.txManager.GetTx(ctx)

	var file fileDataModel.File
	if err := tx.Where("id = ?", id).First(&file).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return fileDomainModel.FileFromModel(ctx, file), nil
}

func (repo *fileRepository) FindListByIds(ctx context.Context, ids []uint64) ([]model.File, error) {
	tx := repo.txManager.GetTx(ctx)

	var files []fileDataModel.File
	if err := tx.Where("id IN ?", ids).Find(&files).Error; err != nil {
		return nil, err
	}

	return fileDomainModel.FileListFromModel(ctx, files), nil
}

func (repo *fileRepository) Save(ctx context.Context, file *model.File) error {
	if file == nil {
		return errors.New("file cannot be nil")
	}

	dataModel := fileDataModel.File{
		MediaType: file.MediaType.ToInt(),
		Name:      file.Name.String(),
		NameIndex: file.NameIndex,
		Path:      file.Path,
		Size:      file.Size,
		Ext:       file.Ext,
		MimeType:  file.MimeType,
		Status:    file.Status,
	}

	tx := repo.txManager.GetTx(ctx)
	if file.ID <= 0 {
		if err := tx.Debug().Create(&dataModel).Error; err != nil {
			return err
		}
		file.ID = dataModel.ID
		file.CreatedAt = dataModel.CreatedAt.Time
	} else {
		// 更新文件信息
		return tx.Debug().Model(&fileDataModel.File{}).Where("id = ?", file.ID).Updates(&dataModel).Error
	}

	return nil
}

func (repo *fileRepository) Delete(ctx context.Context, id uint64) error {
	file, err := repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}

	tx := repo.txManager.GetTx(ctx)

	// 删除文件记录
	return tx.WithContext(ctx).Where("id = ?", id).Delete(&fileDataModel.File{}).Error
}

func (repo *fileRepository) BatchDelete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	tx := repo.txManager.GetTx(ctx)
	// 删除文件记录
	return tx.Where("id IN ?", ids).Delete(&fileDataModel.File{}).Error
}

func (repo *fileRepository) CheckFileNameExists(ctx context.Context, name string, excludeId uint64) (bool, error) {
	tx := repo.txManager.GetTx(ctx)

	query := tx.Debug().Model(&fileDataModel.File{}).Where("name = ?", name)

	// 如果提供了排除ID，则在查询中排除该文件
	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	// 查询是否存在同名文件
	var existingFile fileDataModel.File
	err := query.First(&existingFile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
