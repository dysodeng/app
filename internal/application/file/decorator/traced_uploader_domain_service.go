package decorator

import (
	"context"
	"mime/multipart"

	fileModel "github.com/dysodeng/app/internal/domain/file/model"
	fileDomainSvc "github.com/dysodeng/app/internal/domain/file/service"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

type TracedUploaderDomainService struct {
	inner    fileDomainSvc.UploaderDomainService
	baseSpan string
}

func NewTracedUploaderDomainService(inner fileDomainSvc.UploaderDomainService) fileDomainSvc.UploaderDomainService {
	return &TracedUploaderDomainService{
		inner:    inner,
		baseSpan: "application.file.domain.UploaderDomainService",
	}
}

func (t *TracedUploaderDomainService) UploadFile(ctx context.Context, file *multipart.FileHeader) (*fileModel.File, error) {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".UploadFile")
	defer span.End()
	return t.inner.UploadFile(spanCtx, file)
}

func (t *TracedUploaderDomainService) InitMultipartUpload(ctx context.Context, filename string, fileSize int64) (string, string, error) {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".InitMultipartUpload")
	defer span.End()
	return t.inner.InitMultipartUpload(spanCtx, filename, fileSize)
}

func (t *TracedUploaderDomainService) UploadPart(ctx context.Context, path, uploadId string, partNumber int, file *multipart.FileHeader) (*fileModel.Part, error) {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".UploadPart")
	defer span.End()
	return t.inner.UploadPart(spanCtx, path, uploadId, partNumber, file)
}

func (t *TracedUploaderDomainService) CompleteMultipartUpload(ctx context.Context, uploadId string, parts []fileModel.Part) (*fileModel.File, error) {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".CompleteMultipartUpload")
	defer span.End()
	return t.inner.CompleteMultipartUpload(spanCtx, uploadId, parts)
}

func (t *TracedUploaderDomainService) MultipartUploadStatus(ctx context.Context, uploadId string) ([]fileModel.Part, string, error) {
	spanCtx, span := trace.Tracer().Start(ctx, t.baseSpan+".MultipartUploadStatus")
	defer span.End()
	return t.inner.MultipartUploadStatus(spanCtx, uploadId)
}
