package port

import (
	"context"
	"io"

	"github.com/dysodeng/app/internal/domain/file/model"
)

// FileStorage 文件存储端口
type FileStorage interface {
	TypeByExtension(filePath string) string

	Upload(ctx context.Context, path string, r io.Reader, contentType string) error
	InitMultipartUpload(ctx context.Context, path string, contentType string) (string, error)
	AbortMultipartUpload(ctx context.Context, path, uploadId string) error
	UploadPart(ctx context.Context, path, uploadId string, partNumber int, r io.Reader) (string, error)
	CompleteMultipartUpload(ctx context.Context, path, uploadId string, parts []model.Part) error
	ListUploadedParts(ctx context.Context, path, uploadId string) ([]model.Part, error)

	FullURL(ctx context.Context, path string) string
	RelativePath(ctx context.Context, path string) string
}
