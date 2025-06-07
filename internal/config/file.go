package config

import "fmt"

type ByteSize int64

func (b ByteSize) String() string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%.2fTB", float64(b)/float64(TB))
	case b >= GB:
		return fmt.Sprintf("%.2fGB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2fMB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2fKB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%dB", b)
	}
}

func (b ByteSize) ToInt() int64 {
	return int64(b)
}

const (
	B ByteSize = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)

type FileAllow struct {
	// 允许上传的文件类型
	AllowMimeType []string
	// 允许上传的文件容量大小(单位：字节)
	AllowCapacitySize ByteSize
}

// UserFileAllow 终端用户上传限制
var UserFileAllow = struct {
	Image FileAllow
	Audio FileAllow
	Video FileAllow
}{
	Image: FileAllow{AllowMimeType: []string{"png", "jpg", "jpeg", "gif"}, AllowCapacitySize: 5 * MB},
	Audio: FileAllow{AllowMimeType: []string{"mp3", "wav", "flac", "mid", "mov", "m4a"}, AllowCapacitySize: 20 * MB},
	Video: FileAllow{AllowMimeType: []string{"mp4", "mpg", "avi", "wmv", "mov", "flv", "rmvb", "3gp", "m4v", "mkv"}, AllowCapacitySize: 100 * MB},
}

// AmsFileAllow 管理端上传限制
var AmsFileAllow = struct {
	Image      FileAllow
	Audio      FileAllow
	Video      FileAllow
	Document   FileAllow
	Compressed FileAllow
}{
	Image:      FileAllow{AllowMimeType: []string{"png", "jpg", "jpeg", "gif", "bmp"}, AllowCapacitySize: 5 * MB},
	Audio:      FileAllow{AllowMimeType: []string{"mp3", "wav", "flac", "mid", "mov", "m4a"}, AllowCapacitySize: 20 * MB},
	Video:      FileAllow{AllowMimeType: []string{"mp4", "mpg", "avi", "wmv", "mov", "flv", "rmvb", "3gp", "m4v", "mkv"}, AllowCapacitySize: 2 * GB},
	Document:   FileAllow{AllowMimeType: []string{"pdf", "docx", "doc", "xlsx", "xls", "pptx", "ppt", "cvs", "txt", "markdown", "html"}, AllowCapacitySize: 50 * MB},
	Compressed: FileAllow{AllowMimeType: []string{"zip", "rar", "tar", "tar.gz", "gz", "tgz", "tbz2", "tbz"}, AllowCapacitySize: 500 * MB},
}

// AnonymousFileAllow 匿名用户上传限制
var AnonymousFileAllow = struct {
	File FileAllow
}{
	File: FileAllow{AllowMimeType: []string{"pdf", "doc", "docx", "ppt", "pptx", "xls", "xlsx"}, AllowCapacitySize: 20 * MB},
}
