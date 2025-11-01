package event

import (
	"github.com/google/uuid"

	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
)

// FileUploadedEventType 文件上传事件
const FileUploadedEventType = "file.uploaded"

type FileUploaded struct {
	FileID   uuid.UUID `json:"file_id"`
	FileName string    `json:"file_name"`
	FilePath string    `json:"file_path"`
	FileSize uint64    `json:"file_size"`
}

func NewFileUploadedEvent(fileID uuid.UUID, fileName, filePath string, fileSize uint64) domainEvent.DomainEvent[FileUploaded] {
	payload := FileUploaded{
		FileID:   fileID,
		FileName: fileName,
		FilePath: filePath,
		FileSize: fileSize,
	}
	return domainEvent.NewDomainEvent(FileUploadedEventType, fileID.String(), fileName, payload)
}
