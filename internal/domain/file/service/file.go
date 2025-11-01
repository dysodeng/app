package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/file/errors"
	"github.com/dysodeng/app/internal/domain/file/model"
	"github.com/dysodeng/app/internal/domain/file/repository"
	"github.com/dysodeng/app/internal/domain/file/valueobject"
)

// FileDomainService 文件管理领域服务
type FileDomainService interface {
	// CheckFileNameAvailable 检查文件名是否可用(查重名)
	CheckFileNameAvailable(ctx context.Context, name string, excludeId uuid.UUID) error
	Info(ctx context.Context, id uuid.UUID) (*model.File, error)
	List(ctx context.Context, mediaType valueobject.MediaType, keyword, orderBy, orderType string, page, pageSize int) ([]model.File, int64, error)
	Delete(ctx context.Context, id uuid.UUID, ids []uuid.UUID) error
}

type fileDomainService struct {
	fileRepository repository.FileRepository
}

func NewFileDomainService(fileRepository repository.FileRepository) FileDomainService {
	return &fileDomainService{
		fileRepository: fileRepository,
	}
}

func (svc *fileDomainService) CheckFileNameAvailable(ctx context.Context, name string, excludeId uuid.UUID) error {
	// 检查文件名是否为空
	if name == "" {
		return errors.ErrFileNameEmpty
	}

	// 检查同名文件
	exists, err := svc.fileRepository.CheckFileNameExists(ctx, name, excludeId)
	if err != nil {
		return errors.ErrFileQueryFailed.Wrap(err)
	}
	if exists {
		return errors.ErrFileNameExists
	}

	return nil
}

func (svc *fileDomainService) Info(ctx context.Context, id uuid.UUID) (*model.File, error) {
	if id == uuid.Nil {
		return nil, errors.ErrFileIDEmpty
	}
	file, err := svc.fileRepository.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrFileQueryFailed.Wrap(err)
	}
	return file, nil
}

func (svc *fileDomainService) List(ctx context.Context, mediaType valueobject.MediaType, keyword, orderBy, orderType string, page, pageSize int) ([]model.File, int64, error) {
	// 构建查询参数
	query := repository.FileQuery{
		MediaType: mediaType,
		Keyword:   keyword,
		OrderBy:   orderBy,
		OrderType: orderType,
		Page:      page,
		PageSize:  pageSize,
	}
	list, total, err := svc.fileRepository.FindList(ctx, query)
	if err != nil {
		return nil, 0, errors.ErrFileQueryFailed.Wrap(err)
	}
	return list, total, nil
}

func (svc *fileDomainService) Delete(ctx context.Context, id uuid.UUID, ids []uuid.UUID) error {
	if id != uuid.Nil {
		file, err := svc.fileRepository.FindByID(ctx, id)
		if err != nil {
			return errors.ErrFileQueryFailed.Wrap(err)
		}
		if file.ID == uuid.Nil {
			return errors.ErrFileNotFound
		}
		if err = svc.fileRepository.Delete(ctx, id); err != nil {
			return errors.ErrFileDeleteFailed.Wrap(err)
		}
	}

	if len(ids) > 0 {
		if err := svc.fileRepository.BatchDelete(ctx, ids); err != nil {
			return errors.ErrFileDeleteFailed.Wrap(err)
		}
	}

	return nil
}
