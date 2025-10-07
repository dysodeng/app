package event

import (
	"github.com/dysodeng/app/internal/application/file/event/handler"
)

// HandlerRegistry 事件处理器注册表
type HandlerRegistry struct {
	handlers []any // 类型化事件处理器
}

func NewHandlerRegistry(
	fileUploadedHandler *handler.FileUploadedHandler,
) *HandlerRegistry {
	handlers := make([]any, 0)
	handlers = append(handlers, fileUploadedHandler)
	return &HandlerRegistry{
		handlers: handlers,
	}
}

func (h *HandlerRegistry) Handlers() []any {
	return h.handlers
}
