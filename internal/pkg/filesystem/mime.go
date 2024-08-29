package filesystem

var MimeType = map[string]string{
	"application/xml": "xml",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "xlsx",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "docx",
	"application/vnd.android.package-archive":                                   "apk",
	"application/msword":            "doc",
	"application/ogg":               "ogg",
	"application/pdf":               "pdf",
	"application/vnd.ms-excel":      "xls",
	"application/vnd.ms-powerpoint": "ppt",
	"audio/mpeg":                    "mp2",
	"audio/mp3":                     "mp3",
	"audio/midi":                    "mid",
	"audio/x-wav":                   "wav",
	"audio/x-ms-wma":                "wma",
	"audio/x-ms-wax":                "wax",
	"audio/x-mpegurl":               "m3u",
	"image/bmp":                     "bmp",
	"image/gif":                     "gif",
	"image/ief":                     "ief",
	"image/png":                     "png",
	"image/x-rgb":                   "rgb",
	"image/cgm":                     "cgm",
	"image/x-icon":                  "ico",
	"image/vnd.microsoft.icon":      "ico",
	"image/jp2":                     "jp2",
	"image/jpeg":                    "jpg",
	"image/jpe":                     "jpg",
	"image/jpg":                     "jpg",
	"image/webp":                    "webp",
	"image/svg":                     "svg",
	"image/svg+xml":                 "svg",
	"image/tiff":                    "tif",
	"text/csv":                      "csv",
	"text/plain":                    "txt",
	"text/rtf":                      "rtf",
	"video/quicktime":               "mov",
	"video/mp4":                     "mp4",
	"video/mpeg":                    "mpg",
	"video/x-flv":                   "flv",
	"video/x-ms-wm":                 "wm",
	"video/x-msvideo":               "avi",
	"video/x-sgi-movie":             "movie",
	"video/3gpp":                    "3gp",
}

// IsExistsMimeAllow 判断文件MimeType
func IsExistsMimeAllow(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}
