package model

import (
	"time"

	"github.com/dysodeng/fs"
	"github.com/google/uuid"
)

// MultipartUpload 分片上传信息
type MultipartUpload struct {
	ID        uuid.UUID `json:"id"`
	FileName  string    `json:"file_name"`
	Path      string    `json:"path"`
	Size      uint64    `json:"size"`
	MimeType  string    `json:"mime_type"`
	Ext       string    `json:"ext"`
	UploadID  string    `json:"upload_id"`
	Status    uint8     `json:"status"` // 1-进行中 2-已完成 3-已取消
	Parts     []*Part   `json:"parts"`
	CreatedAt time.Time `json:"created_at"`
}

// Part 分片信息
type Part struct {
	PartNumber int
	ETag       string
	Size       int64
}

// NewMultipartUpload 创建分片上传
func NewMultipartUpload(fileName, path string, size uint64, mimeType, ext, uploadId string) *MultipartUpload {
	return &MultipartUpload{
		FileName: fileName,
		Path:     path,
		Size:     size,
		MimeType: mimeType,
		Ext:      ext,
		UploadID: uploadId,
		Status:   1,
		Parts:    make([]*Part, 0),
	}
}

// AddPart 添加分片
func (m *MultipartUpload) AddPart(partNumber int, etag string, size int64) {
	m.Parts = append(m.Parts, &Part{
		PartNumber: partNumber,
		ETag:       etag,
		Size:       size,
	})
}

// Complete 完成上传
func (m *MultipartUpload) Complete() {
	m.Status = 2
}

// Abort 取消上传
func (m *MultipartUpload) Abort() {
	m.Status = 3
}

func PartListFromStoragePart(parts []fs.MultipartPart) []Part {
	partList := make([]Part, len(parts))
	for i, part := range parts {
		partList[i] = Part{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
			Size:       part.Size,
		}
	}
	return partList
}
