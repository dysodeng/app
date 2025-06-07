package handler

import (
	"context"
	"fmt"
	"log"

	fileEvent "github.com/dysodeng/app/internal/domain/file/event"
	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
)

// FileUploadedHandler 文件上传事件处理器
type FileUploadedHandler struct{}

func NewFileUploadedHandler() *FileUploadedHandler {
	return &FileUploadedHandler{}
}

func (h *FileUploadedHandler) Handle(ctx context.Context, event domainEvent.DomainEvent) error {
	uploadedEvent, ok := event.(*fileEvent.FileUploadedEvent)
	if !ok {
		return fmt.Errorf("invalid event type")
	}
	// 这里可以处理文件上传完成后的额外工作
	log.Printf("file uploaded event id: %s", uploadedEvent.EventID())
	log.Printf("file uploaded event type: %s", uploadedEvent.EventType())
	log.Printf("file uploaded event data: %+v", uploadedEvent.EventData())
	return nil
}

func (h *FileUploadedHandler) CanHandle(eventType string) bool {
	return eventType == fileEvent.FileUploadedEventType
}
