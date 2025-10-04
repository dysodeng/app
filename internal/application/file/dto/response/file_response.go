package response

import (
	"time"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/file/model"
)

// FileResponse 文件响应
type FileResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      uint64    `json:"size"`
	Ext       string    `json:"ext"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

// FromDomainModel 从领域模型转换
func (f *FileResponse) FromDomainModel(file *model.File) {
	f.ID = file.ID
	f.Name = file.Name.String()
	f.Path = file.Path
	f.Size = file.Size
	f.Ext = file.Ext
	f.MimeType = file.MimeType
	f.CreatedAt = file.CreatedAt
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Total int64          `json:"total"`
	Items []FileResponse `json:"items"`
}
