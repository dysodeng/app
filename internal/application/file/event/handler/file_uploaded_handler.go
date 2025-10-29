package handler

import (
	"context"

	fileEvent "github.com/dysodeng/app/internal/domain/file/event"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
)

// FileUploadedHandler 文件上传事件处理器
type FileUploadedHandler struct {
	event.DomainEventHandler[fileEvent.FileUploaded]
}

// NewFileUploadedHandler 创建文件上传事件处理器
func NewFileUploadedHandler() *FileUploadedHandler {
	return &FileUploadedHandler{}
}

// Handle 事件处理
func (h *FileUploadedHandler) Handle(ctx context.Context, event any) error {
	domainEvent, err := h.ParseDomainEvent(ctx, event)
	if err != nil {
		return err
	}

	payload := domainEvent.Payload()

	// 记录事件信息
	logger.Info(ctx, "处理文件上传事件",
		logger.AddField("事件类型", domainEvent.EventType()),
		logger.AddField("发生时间", domainEvent.OccurredAt()),
	)

	// 记录领域聚合根信息
	logger.Info(ctx, "领域事件信息",
		logger.AddField("聚合根ID", domainEvent.AggregateID()),
		logger.AddField("聚合根名称", domainEvent.AggregateName()),
	)

	// 记录文件上传信息
	logger.Info(ctx, "文件已上传",
		logger.AddField("文件ID", payload.FileID.String()),
		logger.AddField("文件名称", payload.FileName),
		logger.AddField("文件路径", payload.FilePath),
		logger.AddField("文件大小", payload.FileSize),
	)

	return nil
}

// InterestedEventTypes 返回感兴趣的事件列表
func (h *FileUploadedHandler) InterestedEventTypes() []string {
	return []string{fileEvent.FileUploadedEventType}
}
