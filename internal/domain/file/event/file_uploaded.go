package event

import (
	sharedEvent "github.com/dysodeng/app/internal/domain/shared/event"
)

const FileUploadedEventType = "file.uploaded"

// FileUploadedEvent 文件上传事件
type FileUploadedEvent struct {
	sharedEvent.BaseDomainEvent
	FileID   uint64 `json:"file_id"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileSize uint64 `json:"file_size"`
}

func NewFileUploadedEvent(fileID uint64, fileName, filePath string, fileSize uint64) *FileUploadedEvent {
	return &FileUploadedEvent{
		BaseDomainEvent: sharedEvent.NewBaseDomainEvent(FileUploadedEventType, fileID),
		FileID:          fileID,
		FileName:        fileName,
		FilePath:        filePath,
		FileSize:        fileSize,
	}
}

func (e *FileUploadedEvent) EventData() map[string]interface{} {
	return map[string]interface{}{
		"file_id":   e.FileID,
		"file_path": e.FilePath,
		"file_name": e.FileName,
		"file_size": e.FileSize,
	}
}
