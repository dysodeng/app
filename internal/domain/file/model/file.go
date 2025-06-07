package model

import (
	"context"
	"strings"
	"time"
	"unicode"

	"github.com/dysodeng/app/internal/infrastructure/persistence/model/file"
	"github.com/dysodeng/app/internal/pkg/model"
	"github.com/dysodeng/app/internal/pkg/storage"
	"github.com/mozillazg/go-pinyin"
)

// MediaType 文件媒体类型
type MediaType uint8

const (
	MediaTypeImage      MediaType = iota + 1 // 图片
	MediaTypeVideo                           // 视频
	MediaTypeAudio                           // 音频
	MediaTypeDocument                        // 文档
	MediaTypeCompressed                      // 压缩文件
)

// ToInt 转换为整数
func (t MediaType) ToInt() uint8 {
	return uint8(t)
}

// String 获取媒体类型描述
func (t MediaType) String() string {
	switch t {
	case MediaTypeImage:
		return "图片"
	case MediaTypeVideo:
		return "视频"
	case MediaTypeAudio:
		return "音频"
	case MediaTypeDocument:
		return "文档"
	case MediaTypeCompressed:
		return "压缩文件"
	default:
		return "未知"
	}
}

// FileName 文件名领域值对象
type FileName string

func (f FileName) String() string {
	return string(f)
}

func (f FileName) Validate() error {
	if f == "" {
		return ErrFileNameEmpty
	}
	return nil
}

// NameIndex 生成文件名索引
func (f FileName) NameIndex() string {
	args := pinyin.Args{Style: pinyin.FirstLetter, Fallback: func(r rune, a pinyin.Args) []string {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return []string{}
		}
		return []string{strings.ToLower(string(r))}
	}}
	list := pinyin.Pinyin(string(f), args)
	if len(list) > 0 {
		var p string
		for _, name := range list {
			p += name[0]
		}
		return p
	}
	return ""
}

// File 文件领域模型
type File struct {
	ID        uint64    `json:"id"`
	MediaType MediaType `json:"media_type"`
	Name      FileName  `json:"name"`
	NameIndex string    `json:"name_index"`
	Path      string    `json:"path"`
	Size      uint64    `json:"size"`
	Ext       string    `json:"ext"`
	MimeType  string    `json:"mime_type"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// NewFile 创建文件
func NewFile(ctx context.Context, name, ext, path, mimeType string, size uint64) *File {
	fileName := FileName(name)
	return &File{
		Name:      fileName,
		NameIndex: fileName.NameIndex(),
		Path:      storage.Instance().RelativePath(ctx, path),
		Ext:       ext,
		MimeType:  mimeType,
		Size:      size,
		Status:    1,
		MediaType: DetermineMediaType(ext, mimeType),
	}
}

func (f *File) Validate() error {
	if f == nil {
		return ErrFileNotFound
	}
	if f.Path == "" {
		return ErrFilePathEmpty
	}
	if err := f.Name.Validate(); err != nil {
		return err
	}
	return nil
}

// DetermineMediaType 根据扩展名和MIME类型确定媒体类型
func DetermineMediaType(ext, mimeType string) MediaType {
	ext = strings.ToLower(strings.TrimLeft(ext, "."))
	// 图片类型
	imageExts := map[string]bool{"jpg": true, "jpeg": true, "png": true, "gif": true, "webp": true}
	if imageExts[strings.ToLower(ext)] || strings.HasPrefix(mimeType, "image/") {
		return MediaTypeImage
	}

	// 视频类型
	videoExts := map[string]bool{"mp4": true, "avi": true, "mov": true, "wmv": true, "flv": true}
	if videoExts[strings.ToLower(ext)] || strings.HasPrefix(mimeType, "video/") {
		return MediaTypeVideo
	}

	// 音频类型
	audioExts := map[string]bool{"mp3": true, "wav": true, "ogg": true, "m4a": true}
	if audioExts[strings.ToLower(ext)] || strings.HasPrefix(mimeType, "audio/") {
		return MediaTypeAudio
	}

	// 压缩包类型
	compressedExts := map[string]bool{"zip": true, "rar": true, "tar": true, "tgz": true, "tar.gz": true, "tbz2": true, "tbz": true, "gz": true}
	if compressedExts[strings.ToLower(ext)] {
		return MediaTypeCompressed
	}

	// 默认为文档类型
	return MediaTypeDocument
}

// ToModel 转换为数据模型
func (f *File) ToModel() *file.File {
	nameIndex := f.NameIndex
	if nameIndex == "" {
		nameIndex = f.Name.NameIndex()
	}
	return &file.File{
		PrimaryKeyID: model.PrimaryKeyID{ID: f.ID},
		MediaType:    f.MediaType.ToInt(),
		Name:         f.Name.String(),
		NameIndex:    nameIndex,
		Path:         f.Path,
		Size:         f.Size,
		Ext:          f.Ext,
		MimeType:     f.MimeType,
		Status:       f.Status,
	}
}

func FileFromModel(ctx context.Context, m file.File) *File {
	return &File{
		ID:        m.ID,
		MediaType: MediaType(m.MediaType),
		Name:      FileName(m.Name),
		NameIndex: m.NameIndex,
		Path:      storage.Instance().FullUrl(ctx, m.Path),
		Size:      m.Size,
		Ext:       m.Ext,
		MimeType:  m.MimeType,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Time,
	}
}

func FileListFromModel(ctx context.Context, files []file.File) []File {
	result := make([]File, len(files))
	for i, m := range files {
		result[i] = File{
			ID:        m.ID,
			MediaType: MediaType(m.MediaType),
			Name:      FileName(m.Name),
			NameIndex: m.NameIndex,
			Path:      storage.Instance().FullUrl(ctx, m.Path),
			Size:      m.Size,
			Ext:       m.Ext,
			MimeType:  m.MimeType,
			Status:    m.Status,
			CreatedAt: m.CreatedAt.Time,
		}
	}
	return result
}
