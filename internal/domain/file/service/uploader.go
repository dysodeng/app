package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/domain/file/event"
	"github.com/dysodeng/app/internal/domain/file/model"
	"github.com/dysodeng/app/internal/domain/file/repository"
	"github.com/dysodeng/app/internal/infrastructure/event/publisher"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/dysodeng/app/internal/pkg/helper"
	fsStorage "github.com/dysodeng/app/internal/pkg/storage"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/fs"
)

// UploaderDomainService 文件上传领域服务
type UploaderDomainService interface {
	// UploadFile 普通文件上传
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*model.File, error)
	// InitMultipartUpload 初始化分片上传
	InitMultipartUpload(ctx context.Context, filename string, fileSize int64) (string, string, error)
	// UploadPart 上传分片
	UploadPart(ctx context.Context, path, uploadId string, partNumber int, file *multipart.FileHeader) (*model.Part, error)
	// CompleteMultipartUpload 完成分片上传
	CompleteMultipartUpload(ctx context.Context, uploadId string, parts []model.Part) (*model.File, error)
	// MultipartUploadStatus 分片上传状态
	MultipartUploadStatus(ctx context.Context, uploadId string) ([]model.Part, string, error)
}

type uploaderDomainService struct {
	baseTraceSpanName  string
	txManager          transactions.TransactionManager
	fileRepository     repository.FileRepository
	uploaderRepository repository.UploaderRepository
	eventPublisher     *publisher.DomainEventPublisher
	storage            *fsStorage.Storage
}

func NewUploaderDomainService(
	txManager transactions.TransactionManager,
	fileRepository repository.FileRepository,
	uploaderRepository repository.UploaderRepository,
	eventPublisher *publisher.DomainEventPublisher,
) UploaderDomainService {
	return &uploaderDomainService{
		baseTraceSpanName:  "domain.file.service.UploaderDomainService",
		txManager:          txManager,
		fileRepository:     fileRepository,
		uploaderRepository: uploaderRepository,
		eventPublisher:     eventPublisher,
		storage:            fsStorage.Instance(),
	}
}

// generateFilePath 生成上传文件路径
func (svc *uploaderDomainService) generateFilePath(ext string) (string, error) {
	if strings.ContainsRune(ext, '/') { // 防止路径注入
		ext = ".invalid"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ext = "." + strings.Replace(ext, ".", "", -1)

	now := time.Now()
	dateDir := now.Format("2006/01/02")

	// 生成唯一文件名部分（纳秒时间戳+随机熵）
	randomBytes := make([]byte, 4) // 4字节提供32位熵值
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// 构建最终文件名（格式：时间戳_随机数.扩展名）
	fileName := fmt.Sprintf("%d_%s%s",
		now.UnixNano(),
		base64.RawURLEncoding.EncodeToString(randomBytes),
		ext,
	)

	return path.Join(
		"resources",
		dateDir,
		fileName,
	), nil
}

// checkFileAllow 检查文件上传限制
func (svc *uploaderDomainService) checkFileAllow(ext, mimeType string, size int64) error {
	ext = strings.TrimLeft(ext, ".")
	mediaType := model.DetermineMediaType(ext, mimeType)

	var allow config.FileAllow
	switch mediaType {
	case model.MediaTypeImage:
		allow = config.AmsFileAllow.Image
	case model.MediaTypeAudio:
		allow = config.AmsFileAllow.Audio
	case model.MediaTypeVideo:
		allow = config.AmsFileAllow.Video
	case model.MediaTypeDocument:
		allow = config.AmsFileAllow.Document
	case model.MediaTypeCompressed:
		allow = config.AmsFileAllow.Compressed
	default:
		return model.ErrFileInvalidType
	}

	if !helper.Contain(allow.AllowMimeType, ext) {
		return model.ErrFileInvalidType
	}
	if size > allow.AllowCapacitySize.ToInt() {
		return model.ErrFileSizeExceeded
	}

	return nil
}

func (svc *uploaderDomainService) UploadFile(ctx context.Context, file *multipart.FileHeader) (*model.File, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".UploadFile")
	defer span.End()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = src.Close()
	}()

	mimeType := fs.TypeByExtension(file.Filename)

	// 检查文件上传限制
	if err = svc.checkFileAllow(ext, mimeType, file.Size); err != nil {
		return nil, err
	}

	// 查询文件是否已存在
	exists, err := svc.fileRepository.CheckFileNameExists(ctx, file.Filename, 0)
	if err != nil {
		return nil, model.ErrFileCheckFailed.Wrap(err)
	}
	if exists {
		return nil, model.ErrFileNameExists
	}

	// 生成最终路径
	filePath, _ := svc.generateFilePath(ext)

	// 文件上传
	uploader := svc.storage.FileSystem().Uploader()
	err = uploader.Upload(spanCtx, filePath, src, fs.WithContentType(mimeType))
	if err != nil {
		return nil, model.ErrFileUploadFailed.Wrap(err)
	}

	// 保存文件记录
	f := model.NewFile(spanCtx, file.Filename, ext, filePath, mimeType, uint64(file.Size))
	if err = f.Validate(); err != nil {
		return nil, err
	}

	if err = svc.fileRepository.Save(spanCtx, f); err != nil {
		return nil, model.ErrFileRecordSaveFailed.Wrap(err)
	}

	f.Path = svc.storage.FullUrl(spanCtx, f.Path)

	// 发布领域事件
	_ = svc.eventPublisher.Publish(spanCtx, event.NewFileUploadedEvent(f.ID, f.Name.String(), f.Path, f.Size))

	return f, nil
}

func (svc *uploaderDomainService) InitMultipartUpload(ctx context.Context, filename string, fileSize int64) (string, string, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".UploadFile")
	defer span.End()

	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := fs.TypeByExtension(filename)

	// 检查文件上传限制
	if err := svc.checkFileAllow(ext, mimeType, fileSize); err != nil {
		return "", "", err
	}

	// 查询文件是否已存在
	exists, err := svc.fileRepository.CheckFileNameExists(ctx, filename, 0)
	if err != nil {
		return "", "", model.ErrFileCheckFailed.Wrap(err)
	}
	if exists {
		return "", "", model.ErrFileNameExists
	}

	filePath, _ := svc.generateFilePath(ext)

	uploader := svc.storage.FileSystem().Uploader()
	uploadId, err := uploader.InitMultipartUpload(spanCtx, filePath, fs.WithContentType(mimeType))
	if err != nil {
		return "", "", model.ErrMultipartInitFailed.Wrap(err)
	}

	mu := model.NewMultipartUpload(filename, filePath, uint64(fileSize), mimeType, ext, uploadId)
	err = svc.uploaderRepository.CreateMultipartUpload(spanCtx, mu)
	if err != nil {
		_ = uploader.AbortMultipartUpload(spanCtx, filePath, uploadId)
		return "", "", model.ErrMultipartInitFailed.Wrap(err)
	}

	return uploadId, svc.storage.FullUrl(spanCtx, filePath), nil
}

func (svc *uploaderDomainService) UploadPart(ctx context.Context, path, uploadId string, partNumber int, file *multipart.FileHeader) (*model.Part, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".UploadPart")
	defer span.End()

	src, err := file.Open()
	if err != nil {
		return nil, model.ErrMultipartReadFailed.Wrap(err)
	}
	defer func() {
		_ = src.Close()
	}()

	uploader := svc.storage.FileSystem().Uploader()
	etag, err := uploader.UploadPart(spanCtx, svc.storage.RelativePath(spanCtx, path), uploadId, partNumber, src)
	if err != nil {
		return nil, model.ErrMultipartUploadFailed.Wrap(err)
	}

	return &model.Part{PartNumber: partNumber, ETag: etag, Size: file.Size}, nil
}

func (svc *uploaderDomainService) CompleteMultipartUpload(ctx context.Context, uploadId string, parts []model.Part) (*model.File, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CompleteMultipartUpload")
	defer span.End()

	mu, err := svc.uploaderRepository.FindMultipartUploadByUploadId(spanCtx, uploadId)
	if err != nil {
		return nil, model.ErrMultipartStatusFailed.Wrap(err)
	}

	var totalSize int64
	var fsParts []fs.MultipartPart
	for _, part := range parts {
		totalSize += part.Size
		fsParts = append(fsParts, fs.MultipartPart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
			Size:       part.Size,
		})
	}

	uploader := svc.storage.FileSystem().Uploader()

	// 检查文件上传限制
	if err = svc.checkFileAllow(mu.Ext, mu.MimeType, totalSize); err != nil {
		_ = uploader.AbortMultipartUpload(spanCtx, mu.Path, uploadId) // 取消分片上传
		_ = svc.uploaderRepository.MultipartUploadStatus(spanCtx, uploadId, 3)
		return nil, err
	}

	filePath := svc.storage.RelativePath(spanCtx, mu.Path)

	if err = uploader.CompleteMultipartUpload(spanCtx, filePath, uploadId, fsParts); err != nil {
		return nil, model.ErrMultipartCompleteFailed.Wrap(err)
	}

	f := model.NewFile(spanCtx, mu.FileName, mu.Ext, filePath, mu.MimeType, uint64(totalSize))
	if err = f.Validate(); err != nil {
		return nil, err
	}

	err = svc.txManager.Transaction(spanCtx, func(txCtx context.Context) error {
		if err = svc.fileRepository.Save(txCtx, f); err != nil {
			return err
		}
		// 处理分片上传状态
		_ = svc.uploaderRepository.MultipartUploadStatus(txCtx, uploadId, 2)
		return nil
	})
	if err != nil {
		_ = uploader.AbortMultipartUpload(spanCtx, filePath, uploadId)
		return nil, model.ErrMultipartCompleteFailed.Wrap(err)
	}

	f.Path = svc.storage.FullUrl(spanCtx, filePath)

	// 发布领域事件
	_ = svc.eventPublisher.Publish(spanCtx, event.NewFileUploadedEvent(f.ID, f.Name.String(), f.Path, f.Size))

	return f, nil
}

func (svc *uploaderDomainService) MultipartUploadStatus(ctx context.Context, uploadId string) ([]model.Part, string, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".MultipartUploadStatus")
	defer span.End()

	mu, err := svc.uploaderRepository.FindMultipartUploadByUploadId(spanCtx, uploadId)
	if err != nil {
		return nil, "", model.ErrMultipartStatusFailed.Wrap(err)
	}

	uploader := svc.storage.FileSystem().Uploader()
	parts, err := uploader.ListUploadedParts(spanCtx, mu.Path, uploadId)
	if err != nil {
		return nil, "", model.ErrMultipartStatusFailed.Wrap(err)
	}

	return model.PartListFromStoragePart(parts), mu.Path, nil
}
