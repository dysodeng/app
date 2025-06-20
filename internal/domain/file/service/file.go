package service

import (
	"context"

	"github.com/dysodeng/app/internal/domain/file/model"
	"github.com/dysodeng/app/internal/domain/file/repository"
	"github.com/dysodeng/app/internal/infrastructure/event/publisher"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
)

// FileDomainService 文件管理领域服务
type FileDomainService interface {
	// CheckFileNameAvailable 检查文件名是否可用(查重名)
	CheckFileNameAvailable(ctx context.Context, name string, excludeId uint64) error
	Info(ctx context.Context, id uint64) (*model.File, error)
	List(ctx context.Context, mediaType model.MediaType, keyword, orderBy, orderType string, page, pageSize int) ([]model.File, int64, error)
	Delete(ctx context.Context, id uint64, ids []uint64) error
}

type fileDomainService struct {
	baseTraceSpanName string
	eventPublisher    *publisher.DomainEventPublisher
	fileRepository    repository.FileRepository
}

func NewFileDomainService(eventPublisher *publisher.DomainEventPublisher, fileRepository repository.FileRepository) FileDomainService {
	return &fileDomainService{
		baseTraceSpanName: "domain.file.service.FileDomainService",
		eventPublisher:    eventPublisher,
		fileRepository:    fileRepository,
	}
}

func (svc *fileDomainService) CheckFileNameAvailable(ctx context.Context, name string, excludeId uint64) error {
	// 检查文件名是否为空
	if name == "" {
		return model.ErrFileNameEmpty
	}

	// 检查同名文件
	exists, err := svc.fileRepository.CheckFileNameExists(ctx, name, excludeId)
	if err != nil {
		return model.ErrFileQueryFailed.Wrap(err)
	}
	if exists {
		return model.ErrFileNameExists
	}

	return nil
}

func (svc *fileDomainService) Info(ctx context.Context, id uint64) (*model.File, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".Info")
	defer span.End()

	if id <= 0 {
		return nil, model.ErrFileIDEmpty
	}
	file, err := svc.fileRepository.FindByID(spanCtx, id)
	if err != nil {
		return nil, model.ErrFileQueryFailed.Wrap(err)
	}
	return file, nil
}

func (svc *fileDomainService) List(ctx context.Context, mediaType model.MediaType, keyword, orderBy, orderType string, page, pageSize int) ([]model.File, int64, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".List")
	defer span.End()

	// 构建查询参数
	query := repository.FileQuery{
		MediaType: mediaType,
		Keyword:   keyword,
		OrderBy:   orderBy,
		OrderType: orderType,
		Page:      page,
		PageSize:  pageSize,
	}

	// 调用仓储接口查询
	list, total, err := svc.fileRepository.FindList(spanCtx, query)
	if err != nil {
		return nil, 0, model.ErrFileQueryFailed.Wrap(err)
	}

	return list, total, nil
}

func (svc *fileDomainService) Delete(ctx context.Context, id uint64, ids []uint64) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".Delete")
	defer span.End()

	if id > 0 {
		file, err := svc.fileRepository.FindByID(spanCtx, id)
		if err != nil {
			return model.ErrFileQueryFailed.Wrap(err)
		}
		if file.ID <= 0 {
			return model.ErrFileNotFound
		}
		if err = svc.fileRepository.Delete(spanCtx, id); err != nil {
			return model.ErrFileDeleteFailed.Wrap(err)
		}
	}

	if len(ids) > 0 {
		if err := svc.fileRepository.BatchDelete(spanCtx, ids); err != nil {
			return model.ErrFileDeleteFailed.Wrap(err)
		}
	}

	return nil
}
