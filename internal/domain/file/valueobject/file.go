package valueobject

import (
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"

	"github.com/dysodeng/app/internal/domain/file/errors"
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
		return errors.ErrFileNameEmpty
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
