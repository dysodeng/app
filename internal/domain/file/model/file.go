package model

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/file/errors"
	"github.com/dysodeng/app/internal/domain/file/valueobject"
)

// File 文件领域模型
type File struct {
	ID        uuid.UUID             `json:"id"`
	MediaType valueobject.MediaType `json:"media_type"`
	Name      valueobject.FileName  `json:"name"`
	NameIndex string                `json:"name_index"`
	Path      string                `json:"path"`
	Size      uint64                `json:"size"`
	Ext       string                `json:"ext"`
	MimeType  string                `json:"mime_type"`
	Status    uint8                 `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
}

// NewFile 创建文件（不再依赖基础设施进行路径转换）
func NewFile(name, ext, path, mimeType string, size uint64) *File {
	fileName := valueobject.FileName(name)
	return &File{
		Name:      fileName,
		NameIndex: fileName.NameIndex(),
		Path:      path,
		Ext:       ext,
		MimeType:  mimeType,
		Size:      size,
		Status:    1,
		MediaType: DetermineMediaType(ext, mimeType),
	}
}

func (f *File) Validate() error {
	if f == nil {
		return errors.ErrFileNotFound
	}
	if f.Path == "" {
		return errors.ErrFilePathEmpty
	}
	if err := f.Name.Validate(); err != nil {
		return err
	}
	return nil
}

// DetermineMediaType 根据扩展名和MIME类型确定媒体类型
func DetermineMediaType(ext, mimeType string) valueobject.MediaType {
	ext = strings.ToLower(strings.TrimLeft(ext, "."))
	// 图片类型
	imageExts := map[string]bool{"jpg": true, "jpeg": true, "png": true, "gif": true, "webp": true}
	if imageExts[strings.ToLower(ext)] || strings.HasPrefix(mimeType, "image/") {
		return valueobject.MediaTypeImage
	}

	// 视频类型
	videoExts := map[string]bool{"mp4": true, "avi": true, "mov": true, "wmv": true, "flv": true}
	if videoExts[strings.ToLower(ext)] || strings.HasPrefix(mimeType, "video/") {
		return valueobject.MediaTypeVideo
	}

	// 音频类型
	audioExts := map[string]bool{"mp3": true, "wav": true, "ogg": true, "m4a": true}
	if audioExts[strings.ToLower(ext)] || strings.HasPrefix(mimeType, "audio/") {
		return valueobject.MediaTypeAudio
	}

	// 压缩包类型
	compressedExts := map[string]bool{"zip": true, "rar": true, "tar": true, "tgz": true, "tar.gz": true, "tbz2": true, "tbz": true, "gz": true}
	if compressedExts[strings.ToLower(ext)] {
		return valueobject.MediaTypeCompressed
	}

	// 默认为文档类型
	return valueobject.MediaTypeDocument
}
