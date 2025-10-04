package file

import (
	"github.com/dysodeng/app/internal/infrastructure/shared/model"
)

// File 文件记录
type File struct {
	model.DistributedPrimaryKeyID
	MediaType uint8  `gorm:"not null;default:0;comment:媒体类型 1-图片 2-视频 3-音频 4-文档 5-压缩文件" json:"media_type"`
	NameIndex string `gorm:"type:varchar(150);index:file_name_index_idx;not null;default:'';comment:文件名称索引" json:"name_index"`
	Name      string `gorm:"type:varchar(150);index:file_name_idx;not null;default:'';comment:文件名称" json:"name"`
	Path      string `gorm:"type:varchar(255);not null;default:'';comment:文件路径" json:"path"`
	Size      uint64 `gorm:"type:bigint;not null;default:0;comment:文件大小(字节)" json:"size"`
	Ext       string `gorm:"type:varchar(10);not null;default:'';comment:文件扩展名" json:"ext"`
	MimeType  string `gorm:"type:varchar(50);not null;default:'';comment:文件MIME类型" json:"mime_type"`
	Status    uint8  `gorm:"not null;default:1;comment:文件状态 1-正常" json:"status"`
	model.Time
}

func (File) TableName() string {
	return "files"
}

// MultipartUpload 文件分片上传记录
type MultipartUpload struct {
	model.DistributedPrimaryKeyID
	UploadID string `gorm:"type:varchar(100);index:file_mu_idx,unique;not null;default:'';comment:上传ID" json:"upload_id"`
	FileName string `gorm:"type:varchar(255);not null;default:'';comment:文件名" json:"file_name"`
	Path     string `gorm:"type:varchar(255);not null;comment:文件存储路径" json:"path"`
	Size     uint64 `gorm:"type:bigint;not null;default:0;comment:文件大小(字节)" json:"size"`
	MimeType string `gorm:"type:varchar(50);not null;default:'';comment:mime类型" json:"mime_type"`
	Ext      string `gorm:"type:varchar(10);not null;default:'';comment:文件扩展名" json:"ext"`
	Status   uint8  `gorm:"not null;default:1;comment:状态 1-进行中 2-已完成 3-已取消" json:"status"`
	model.Time
}

func (MultipartUpload) TableName() string {
	return "files_multipart_uploads"
}
