package config

type FileAllow struct {
	// 允许上传的文件类型
	AllowMimeType []string
	// 允许上传的文件容量大小(单位：字节)
	AllowCapacitySize int64
}

const (
	FiveMB    int64 = 5 * 1024 * 1024   // 5MB
	TwentyMB  int64 = 20 * 1024 * 1024  // 20MB
	HundredMB int64 = 100 * 1024 * 1024 // 100MB
)

// UserFileAllow 终端用户上传限制
var UserFileAllow = struct {
	Image FileAllow
	Audio FileAllow
	Video FileAllow
}{
	Image: FileAllow{AllowMimeType: []string{"png", "jpg", "jpeg", "gif"}, AllowCapacitySize: FiveMB},
	Audio: FileAllow{AllowMimeType: []string{"mp3", "wav", "flac", "mid", "mov", "m4a"}, AllowCapacitySize: TwentyMB},
	Video: FileAllow{AllowMimeType: []string{"mp4", "mpg", "avi", "wmv", "mov", "flv", "rmvb", "3gp", "m4v", "mkv"}, AllowCapacitySize: HundredMB},
}

// AmsFileAllow 管理端上传限制
var AmsFileAllow = struct {
	Image      FileAllow
	Audio      FileAllow
	Video      FileAllow
	RewardFile FileAllow
}{
	Image:      FileAllow{AllowMimeType: []string{"png", "jpg", "jpeg", "gif", "bmp"}, AllowCapacitySize: FiveMB},
	Audio:      FileAllow{AllowMimeType: []string{"mp3", "wav", "flac", "mid", "mov", "m4a"}, AllowCapacitySize: TwentyMB},
	Video:      FileAllow{AllowMimeType: []string{"mp4", "mpg", "avi", "wmv", "mov", "flv", "rmvb", "3gp", "m4v", "mkv"}, AllowCapacitySize: HundredMB},
	RewardFile: FileAllow{AllowMimeType: []string{"xlsx", "xls"}, AllowCapacitySize: FiveMB},
}

// AnonymousFileAllow 匿名用户上传限制
var AnonymousFileAllow = struct {
	File FileAllow
}{
	File: FileAllow{AllowMimeType: []string{"pdf", "doc", "docx", "ppt", "pptx", "xls", "xlsx"}, AllowCapacitySize: TwentyMB},
}
