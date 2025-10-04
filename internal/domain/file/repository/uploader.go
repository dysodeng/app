package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/file/model"
)

// UploaderRepository 文件上传仓储接口
type UploaderRepository interface {
	// CreateMultipartUpload 创建分片上传记录
	CreateMultipartUpload(ctx context.Context, mu *model.MultipartUpload) error
	// FindMultipartUploadByUploadId 根据分片上传id查询上传记录
	FindMultipartUploadByUploadId(ctx context.Context, uploadId string) (*model.MultipartUpload, error)
	// MultipartUploadStatus 分片上传状态设置
	MultipartUploadStatus(ctx context.Context, uploadId string, status uint8) error
}
