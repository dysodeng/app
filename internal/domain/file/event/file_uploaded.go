package event

import (
	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/infrastructure/event"
)

// FileUploadedEventType 文件上传事件
const FileUploadedEventType = "file.uploaded"

type FileUploaded struct {
	FileID   uuid.UUID `json:"file_id"`
	FileName string    `json:"file_name"`
	FilePath string    `json:"file_path"`
	FileSize uint64    `json:"file_size"`
}

func NewFileUploadedEvent(fileID uuid.UUID, fileName, filePath string, fileSize uint64) event.DomainEvent[FileUploaded] {
	// 创建FileUploaded数据
	payload := FileUploaded{
		FileID:   fileID,
		FileName: fileName,
		FilePath: filePath,
		FileSize: fileSize,
	}

	return event.NewDomainEvent(FileUploadedEventType, fileID.String(), fileName, payload)
}
