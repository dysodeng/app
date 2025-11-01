package file

import (
	"context"
	"io"

	"github.com/dysodeng/fs"

	domainModel "github.com/dysodeng/app/internal/domain/file/model"
	domainPort "github.com/dysodeng/app/internal/domain/file/port"
	infraStorage "github.com/dysodeng/app/internal/infrastructure/shared/storage"
)

// StorageAdapter 文件存储端口适配器
type StorageAdapter struct {
	st *infraStorage.Storage
}

func NewFileStorageAdapter(st *infraStorage.Storage) domainPort.FileStorage {
	return &StorageAdapter{st: st}
}

func (adapter *StorageAdapter) TypeByExtension(filePath string) string {
	return fs.TypeByExtension(filePath)
}

func (adapter *StorageAdapter) Upload(ctx context.Context, path string, r io.Reader, contentType string) error {
	return adapter.st.FileSystem().Uploader().Upload(ctx, path, r, fs.WithContentType(contentType))
}

func (adapter *StorageAdapter) InitMultipartUpload(ctx context.Context, path string, contentType string) (string, error) {
	return adapter.st.FileSystem().Uploader().InitMultipartUpload(ctx, path, fs.WithContentType(contentType))
}

func (adapter *StorageAdapter) AbortMultipartUpload(ctx context.Context, path, uploadId string) error {
	return adapter.st.FileSystem().Uploader().AbortMultipartUpload(ctx, path, uploadId)
}

func (adapter *StorageAdapter) UploadPart(ctx context.Context, path, uploadId string, partNumber int, r io.Reader) (string, error) {
	return adapter.st.FileSystem().Uploader().UploadPart(ctx, path, uploadId, partNumber, r)
}

func (adapter *StorageAdapter) CompleteMultipartUpload(ctx context.Context, path, uploadId string, parts []domainModel.Part) error {
	// 适配 domain.Part 到 fs.MultipartPart
	fsParts := make([]fs.MultipartPart, 0, len(parts))
	for _, p := range parts {
		fsParts = append(fsParts, fs.MultipartPart{PartNumber: p.PartNumber, ETag: p.ETag, Size: p.Size})
	}
	return adapter.st.FileSystem().Uploader().CompleteMultipartUpload(ctx, path, uploadId, fsParts)
}

func (adapter *StorageAdapter) ListUploadedParts(ctx context.Context, path, uploadId string) ([]domainModel.Part, error) {
	fsParts, err := adapter.st.FileSystem().Uploader().ListUploadedParts(ctx, path, uploadId)
	if err != nil {
		return nil, err
	}

	parts := make([]domainModel.Part, len(fsParts))
	for i, p := range fsParts {
		parts[i] = domainModel.Part{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
			Size:       p.Size,
		}
	}

	return parts, nil
}

func (adapter *StorageAdapter) FullURL(ctx context.Context, path string) string {
	return adapter.st.FullUrl(ctx, path)
}

func (adapter *StorageAdapter) RelativePath(ctx context.Context, path string) string {
	return adapter.st.RelativePath(ctx, path)
}
