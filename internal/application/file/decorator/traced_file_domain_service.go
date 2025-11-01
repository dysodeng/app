package decorator

import (
	"context"

	"github.com/google/uuid"

	fileModel "github.com/dysodeng/app/internal/domain/file/model"
	fileDomainSvc "github.com/dysodeng/app/internal/domain/file/service"
	fileVO "github.com/dysodeng/app/internal/domain/file/valueobject"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

type TracedFileDomainService struct {
	inner    fileDomainSvc.FileDomainService
	baseSpan string
}

func NewTracedFileDomainService(inner fileDomainSvc.FileDomainService) fileDomainSvc.FileDomainService {
	return &TracedFileDomainService{
		inner:    inner,
		baseSpan: "application.file.domain.FileDomainService",
	}
}

func (t *TracedFileDomainService) CheckFileNameAvailable(ctx context.Context, name string, excludeId uuid.UUID) error {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".CheckFileNameAvailable")
	defer span.End()
	return t.inner.CheckFileNameAvailable(spanCtx, name, excludeId)
}

func (t *TracedFileDomainService) Info(ctx context.Context, id uuid.UUID) (*fileModel.File, error) {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".Info")
	defer span.End()
	return t.inner.Info(spanCtx, id)
}

func (t *TracedFileDomainService) List(ctx context.Context, mediaType fileVO.MediaType, keyword, orderBy, orderType string, page, pageSize int) ([]fileModel.File, int64, error) {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".List")
	defer span.End()
	return t.inner.List(spanCtx, mediaType, keyword, orderBy, orderType, page, pageSize)
}

func (t *TracedFileDomainService) Delete(ctx context.Context, id uuid.UUID, ids []uuid.UUID) error {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".Delete")
	defer span.End()
	return t.inner.Delete(spanCtx, id, ids)
}
