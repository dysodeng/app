package service

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/dysodeng/fs"

	"github.com/dysodeng/app/internal/application/file/dto/command"
	"github.com/dysodeng/app/internal/application/file/dto/response"
	fileErrors "github.com/dysodeng/app/internal/domain/file/errors"
	fileEvent "github.com/dysodeng/app/internal/domain/file/event"
	fileModel "github.com/dysodeng/app/internal/domain/file/model"
	filePort "github.com/dysodeng/app/internal/domain/file/port"
	fileRepository "github.com/dysodeng/app/internal/domain/file/repository"
	"github.com/dysodeng/app/internal/domain/file/service"
	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
	sharedPort "github.com/dysodeng/app/internal/domain/shared/port"
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

// uploaderApplicationService 结构体
type uploaderApplicationService struct {
	baseTraceSpanName  string
	uploaderService    service.UploaderDomainService
	eventPublisher     sharedPort.EventPublisher
	txManager          sharedPort.TransactionManager
	fileRepository     fileRepository.FileRepository
	uploaderRepository fileRepository.UploaderRepository
	storage            filePort.FileStorage
}

func NewUploaderApplicationService(
	uploaderService service.UploaderDomainService,
	eventPublisher sharedPort.EventPublisher,
	txManager sharedPort.TransactionManager,
	fileRepository fileRepository.FileRepository,
	uploaderRepository fileRepository.UploaderRepository,
	storage filePort.FileStorage,
) UploaderApplicationService {
	return &uploaderApplicationService{
		baseTraceSpanName:  "application.file.UploaderApplicationService",
		uploaderService:    uploaderService,
		eventPublisher:     eventPublisher,
		txManager:          txManager,
		fileRepository:     fileRepository,
		uploaderRepository: uploaderRepository,
		storage:            storage,
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

	// 持久化（事务）
	if err = svc.txManager.Transaction(spanCtx, func(txCtx context.Context) error {
		return svc.fileRepository.Save(txCtx, f)
	}); err != nil {
		logger.Error(spanCtx, "保存文件记录失败", logger.ErrorField(err))
		return nil, fileErrors.ErrFileRecordSaveFailed.Wrap(err)
	}

	f.Path = svc.storage.FullURL(spanCtx, f.Path)

	// 发布领域事件
	evt := fileEvent.NewFileUploadedEvent(f.ID, f.Name.String(), f.Path, f.Size)
	if err = svc.eventPublisher.Publish(spanCtx, domainEvent.DomainEvent[any]{
		Type:          evt.Type,
		AggregateID:   evt.AggregateID,
		AggregateName: evt.AggregateName,
		Payload:       evt.Payload,
	}); err != nil {
		logger.Warn(spanCtx, "发布文件上传事件失败", logger.ErrorField(err))
	}

	fileRes := &response.FileResponse{}
	fileRes.FromDomainModel(f)

	return fileRes, nil
}

func (svc *uploaderApplicationService) InitMultipartUpload(ctx context.Context, filename string, fileSize int64) (*response.InitMultipartUploadResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".InitMultipartUpload")
	defer span.End()

	uploadId, relPath, err := svc.uploaderService.InitMultipartUpload(spanCtx, filename, fileSize)
	if err != nil {
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := fs.TypeByExtension(filename)
	mu := fileModel.NewMultipartUpload(filename, relPath, uint64(fileSize), mimeType, ext, uploadId)

	if err = svc.uploaderRepository.CreateMultipartUpload(spanCtx, mu); err != nil {
		_ = svc.storage.AbortMultipartUpload(spanCtx, relPath, uploadId)
		logger.Error(spanCtx, "创建分片上传记录失败", logger.ErrorField(err))
		return nil, fileErrors.ErrMultipartInitFailed.Wrap(err)
	}

	return &response.InitMultipartUploadResponse{
		UploadId: uploadId,
		Path:     svc.storage.FullURL(spanCtx, relPath),
	}, nil
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
		// 领域校验或存储合并失败：尝试回滚并置取消
		mu, _ := svc.uploaderRepository.FindMultipartUploadByUploadId(spanCtx, uploadId)
		if mu != nil {
			_ = svc.storage.AbortMultipartUpload(spanCtx, mu.Path, uploadId)
			_ = svc.uploaderRepository.MultipartUploadStatus(spanCtx, uploadId, 3)
		}
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	// 持久化文件记录与更新状态
	if err = svc.txManager.Transaction(spanCtx, func(txCtx context.Context) error {
		if err := svc.fileRepository.Save(txCtx, f); err != nil {
			return err
		}
		return svc.uploaderRepository.MultipartUploadStatus(txCtx, uploadId, 2)
	}); err != nil {
		_ = svc.storage.AbortMultipartUpload(spanCtx, f.Path, uploadId)
		_ = svc.uploaderRepository.MultipartUploadStatus(spanCtx, uploadId, 3)
		logger.Error(spanCtx, "完成分片上传持久化失败", logger.ErrorField(err))
		return nil, fileErrors.ErrMultipartCompleteFailed.Wrap(err)
	}

	f.Path = svc.storage.FullURL(spanCtx, f.Path)

	evt := fileEvent.NewFileUploadedEvent(f.ID, f.Name.String(), f.Path, f.Size)
	if err = svc.eventPublisher.Publish(spanCtx, domainEvent.DomainEvent[any]{
		Type:          evt.Type,
		AggregateID:   evt.AggregateID,
		AggregateName: evt.AggregateName,
		Payload:       evt.Payload,
	}); err != nil {
		logger.Warn(spanCtx, "发布文件上传事件失败", logger.ErrorField(err))
	}

	fileRes := &response.FileResponse{}
	fileRes.FromDomainModel(f)
	return fileRes, nil
}

func (svc *uploaderApplicationService) MultipartUploadStatus(ctx context.Context, uploadId string) (*response.MultipartUploadStatusResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".MultipartUploadStatus")
	defer span.End()

	parts, relPath, err := svc.uploaderService.MultipartUploadStatus(spanCtx, uploadId)
	if err != nil {
		logger.Error(spanCtx, err.Error(), logger.ErrorField(err))
		return nil, err
	}

	return &response.MultipartUploadStatusResponse{
		Parts: response.PartListFormDomainModel(parts),
		Path:  svc.storage.FullURL(spanCtx, relPath),
	}, nil
}
