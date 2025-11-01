package port

import "github.com/dysodeng/app/internal/domain/file/valueobject"

// FilePolicy 文件上传策略端口
type FilePolicy interface {
	// Allow 返回指定媒体类型的允许后缀与容量限制
	Allow(mediaType valueobject.MediaType) (allowedExts []string, maxSize int64)
}
