package file

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/dysodeng/app/internal/domain/file/model"
	fileDomainRepository "github.com/dysodeng/app/internal/domain/file/repository"
	"github.com/dysodeng/app/internal/domain/file/valueobject"
	"github.com/dysodeng/app/internal/infrastructure/persistence/entity/file"
	"github.com/dysodeng/app/internal/infrastructure/persistence/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
)

type fileRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewFileRepository(txManager transactions.TransactionManager) fileDomainRepository.FileRepository {
	return &fileRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.file.FileRepository",
		txManager:         txManager,
	}
}

func (repo *fileRepository) FindList(ctx context.Context, query fileDomainRepository.FileQuery) ([]model.File, int64, error) {
	tx := repo.txManager.GetTx(ctx)

	// 构建查询条件
	db := tx.Debug().Model(&file.File{})

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
	var files []file.File
	if err := db.Find(&files).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	return repo.fileListFromModel(ctx, files), total, nil
}

func (repo *fileRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.File, error) {
	tx := repo.txManager.GetTx(ctx)

	var f file.File
	if err := tx.Where("id = ?", id).First(&f).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return repo.fileFromModel(ctx, f), nil
}

func (repo *fileRepository) FindListByIds(ctx context.Context, ids []uuid.UUID) ([]model.File, error) {
	tx := repo.txManager.GetTx(ctx)

	var files []file.File
	if err := tx.Where("id IN ?", ids).Find(&files).Error; err != nil {
		return nil, err
	}

	return repo.fileListFromModel(ctx, files), nil
}

func (repo *fileRepository) Save(ctx context.Context, f *model.File) error {
	if f == nil {
		return errors.New("file cannot be nil")
	}

	dataModel := file.File{
		MediaType: f.MediaType.ToInt(),
		Name:      f.Name.String(),
		NameIndex: f.NameIndex,
		Path:      f.Path,
		Size:      f.Size,
		Ext:       f.Ext,
		MimeType:  f.MimeType,
		Status:    f.Status,
	}

	tx := repo.txManager.GetTx(ctx)
	if f.ID == uuid.Nil {
		if err := tx.Debug().Create(&dataModel).Error; err != nil {
			return err
		}
		f.ID = dataModel.ID
		f.CreatedAt = dataModel.CreatedAt.Time
	} else {
		// 更新文件信息
		return tx.Debug().Model(&file.File{}).Where("id = ?", f.ID).Updates(&dataModel).Error
	}

	return nil
}

func (repo *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	f, err := repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if f == nil {
		return nil
	}

	tx := repo.txManager.GetTx(ctx)

	// 删除文件记录
	return tx.WithContext(ctx).Where("id = ?", id).Delete(&file.File{}).Error
}

func (repo *fileRepository) BatchDelete(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	tx := repo.txManager.GetTx(ctx)
	// 删除文件记录
	return tx.Where("id IN ?", ids).Delete(&file.File{}).Error
}

func (repo *fileRepository) CheckFileNameExists(ctx context.Context, name string, excludeId uuid.UUID) (bool, error) {
	return false, nil
}

func (repo *fileRepository) fileFromModel(ctx context.Context, m file.File) *model.File {
	return &model.File{
		ID:        m.ID,
		MediaType: valueobject.MediaType(m.MediaType),
		Name:      valueobject.FileName(m.Name),
		NameIndex: m.NameIndex,
		Path:      storage.Instance().FullUrl(ctx, m.Path),
		Size:      m.Size,
		Ext:       m.Ext,
		MimeType:  m.MimeType,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Time,
	}
}

func (repo *fileRepository) fileListFromModel(ctx context.Context, files []file.File) []model.File {
	result := make([]model.File, len(files))
	for i, m := range files {
		result[i] = *repo.fileFromModel(ctx, m)
	}
	return result
}
