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

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/file/errors"
	"github.com/dysodeng/app/internal/domain/file/model"
	filePort "github.com/dysodeng/app/internal/domain/file/port"
	"github.com/dysodeng/app/internal/domain/file/repository"
	"github.com/dysodeng/app/internal/infrastructure/shared/helper"
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
	fileRepository     repository.FileRepository
	uploaderRepository repository.UploaderRepository
	storage            filePort.FileStorage
	policy             filePort.FilePolicy
}

func NewUploaderDomainService(
	fileRepository repository.FileRepository,
	uploaderRepository repository.UploaderRepository,
	storage filePort.FileStorage,
	policy filePort.FilePolicy,
) UploaderDomainService {
	return &uploaderDomainService{
		fileRepository:     fileRepository,
		uploaderRepository: uploaderRepository,
		storage:            storage,
		policy:             policy,
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
	ext = "." + strings.ReplaceAll(ext, ".", "")

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

	allowedExts, maxSize := svc.policy.Allow(mediaType)
	if !helper.Contain(allowedExts, ext) {
		return errors.ErrFileInvalidType
	}
	if size > maxSize {
		return errors.ErrFileSizeExceeded
	}
	return nil
}

func (svc *uploaderDomainService) UploadFile(ctx context.Context, file *multipart.FileHeader) (*model.File, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = src.Close() }()

	mimeType := svc.storage.TypeByExtension(file.Filename)

	// 检查文件上传限制
	if err = svc.checkFileAllow(ext, mimeType, file.Size); err != nil {
		return nil, err
	}

	// 查重
	exists, err := svc.fileRepository.CheckFileNameExists(ctx, file.Filename, uuid.Nil)
	if err != nil {
		return nil, errors.ErrFileCheckFailed.Wrap(err)
	}
	if exists {
		return nil, errors.ErrFileNameExists
	}

	// 生成最终路径（相对路径）
	filePath, _ := svc.generateFilePath(ext)

	// 上传
	if err = svc.storage.Upload(ctx, filePath, src, mimeType); err != nil {
		return nil, errors.ErrFileUploadFailed.Wrap(err)
	}

	f := model.NewFile(file.Filename, ext, filePath, mimeType, uint64(file.Size))
	if err = f.Validate(); err != nil {
		return nil, err
	}
	return f, nil
}

func (svc *uploaderDomainService) InitMultipartUpload(ctx context.Context, filename string, fileSize int64) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := svc.storage.TypeByExtension(filename)

	if err := svc.checkFileAllow(ext, mimeType, fileSize); err != nil {
		return "", "", err
	}

	// 查重
	exists, err := svc.fileRepository.CheckFileNameExists(ctx, filename, uuid.Nil)
	if err != nil {
		return "", "", errors.ErrFileCheckFailed.Wrap(err)
	}
	if exists {
		return "", "", errors.ErrFileNameExists
	}

	filePath, _ := svc.generateFilePath(ext)

	uploadId, err := svc.storage.InitMultipartUpload(ctx, filePath, mimeType)
	if err != nil {
		return "", "", errors.ErrMultipartInitFailed.Wrap(err)
	}

	return uploadId, filePath, nil
}

func (svc *uploaderDomainService) UploadPart(ctx context.Context, path, uploadId string, partNumber int, file *multipart.FileHeader) (*model.Part, error) {
	src, err := file.Open()
	if err != nil {
		return nil, errors.ErrMultipartReadFailed.Wrap(err)
	}
	defer func() { _ = src.Close() }()

	etag, err := svc.storage.UploadPart(ctx, svc.storage.RelativePath(ctx, path), uploadId, partNumber, src)
	if err != nil {
		return nil, errors.ErrMultipartUploadFailed.Wrap(err)
	}
	return &model.Part{PartNumber: partNumber, ETag: etag, Size: file.Size}, nil
}

func (svc *uploaderDomainService) CompleteMultipartUpload(ctx context.Context, uploadId string, parts []model.Part) (*model.File, error) {
	mu, err := svc.uploaderRepository.FindMultipartUploadByUploadId(ctx, uploadId)
	if err != nil {
		return nil, errors.ErrMultipartStatusFailed.Wrap(err)
	}

	var totalSize int64
	for _, p := range parts {
		totalSize += p.Size
	}

	// 检查限制
	if err = svc.checkFileAllow(mu.Ext, mu.MimeType, totalSize); err != nil {
		return nil, err
	}

	filePath := svc.storage.RelativePath(ctx, mu.Path)
	if err = svc.storage.CompleteMultipartUpload(ctx, filePath, uploadId, parts); err != nil {
		return nil, errors.ErrMultipartCompleteFailed.Wrap(err)
	}

	f := model.NewFile(mu.FileName, mu.Ext, filePath, mu.MimeType, uint64(totalSize))
	if err = f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

func (svc *uploaderDomainService) MultipartUploadStatus(ctx context.Context, uploadId string) ([]model.Part, string, error) {
	mu, err := svc.uploaderRepository.FindMultipartUploadByUploadId(ctx, uploadId)
	if err != nil {
		return nil, "", errors.ErrMultipartStatusFailed.Wrap(err)
	}

	parts, err := svc.storage.ListUploadedParts(ctx, mu.Path, uploadId)
	if err != nil {
		return nil, "", errors.ErrMultipartStatusFailed.Wrap(err)
	}

	return parts, mu.Path, nil
}
