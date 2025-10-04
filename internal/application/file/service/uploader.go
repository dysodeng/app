package service

import (
	"context"
	"mime/multipart"

	"github.com/dysodeng/app/internal/application/file/dto/command"
	"github.com/dysodeng/app/internal/application/file/dto/response"
	"github.com/dysodeng/app/internal/domain/file/service"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

// UploaderApplicationService 文件上传应用服务
type UploaderApplicationService interface {
	// UploadFile 上传文件
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*response.FileResponse, error)
	// InitMultipartUpload 初始化分片上传
	InitMultipartUpload(ctx context.Context, filename string, fileSize int64) (*response.InitMultipartUploadResponse, error)
	// UploadPart 上传分片
	UploadPart(ctx context.Context, path, uploadId string, partNumber int, fileHeader *multipart.FileHeader) (*response.Part, error)
	// CompleteMultipartUpload 完成分片上传
	CompleteMultipartUpload(ctx context.Context, uploadId string, parts []command.Part) (*response.FileResponse, error)
	// MultipartUploadStatus 查询分片上传状态
	MultipartUploadStatus(ctx context.Context, uploadId string) (*response.MultipartUploadStatusResponse, error)
}

type uploaderApplicationService struct {
	baseTraceSpanName string
	uploaderService   service.UploaderDomainService
}

func NewUploaderApplicationService(uploaderService service.UploaderDomainService) UploaderApplicationService {
	return &uploaderApplicationService{
		baseTraceSpanName: "application.file.UploaderApplicationService",
		uploaderService:   uploaderService,
	}
}

func (svc *uploaderApplicationService) UploadFile(ctx context.Context, file *multipart.FileHeader) (*response.FileResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".UploadFile")
	defer span.End()

	f, err := svc.uploaderService.UploadFile(spanCtx, file)
	if err != nil {
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	fileRes := &response.FileResponse{}
	fileRes.FromDomainModel(f)
	return fileRes, nil
}

func (svc *uploaderApplicationService) InitMultipartUpload(ctx context.Context, filename string, fileSize int64) (*response.InitMultipartUploadResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".InitMultipartUpload")
	defer span.End()

	uploadId, path, err := svc.uploaderService.InitMultipartUpload(spanCtx, filename, fileSize)
	if err != nil {
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	return &response.InitMultipartUploadResponse{UploadId: uploadId, Path: path}, nil
}

func (svc *uploaderApplicationService) UploadPart(ctx context.Context, path, uploadId string, partNumber int, fileHeader *multipart.FileHeader) (*response.Part, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".UploadPart")
	defer span.End()

	part, err := svc.uploaderService.UploadPart(spanCtx, path, uploadId, partNumber, fileHeader)
	if err != nil {
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	return &response.Part{
		PartNumber: partNumber,
		ETag:       part.ETag,
		Size:       part.Size,
	}, nil
}

func (svc *uploaderApplicationService) CompleteMultipartUpload(ctx context.Context, uploadId string, parts []command.Part) (*response.FileResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CompleteMultipartUpload")
	defer span.End()

	f, err := svc.uploaderService.CompleteMultipartUpload(spanCtx, uploadId, command.PartList(parts).ToDomainModel())
	if err != nil {
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	fileRes := &response.FileResponse{}
	fileRes.FromDomainModel(f)
	return fileRes, nil
}

func (svc *uploaderApplicationService) MultipartUploadStatus(ctx context.Context, uploadId string) (*response.MultipartUploadStatusResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".MultipartUploadStatus")
	defer span.End()

	parts, path, err := svc.uploaderService.MultipartUploadStatus(spanCtx, uploadId)
	if err != nil {
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	return &response.MultipartUploadStatusResponse{
		Parts: response.PartListFormDomainModel(parts),
		Path:  path,
	}, nil
}
