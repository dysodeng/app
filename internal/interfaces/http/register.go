package http

import (
	"github.com/dysodeng/app/internal/interfaces/http/handler"
	"github.com/dysodeng/app/internal/interfaces/http/handler/file"
)

// HandlerRegistry 控制器注册表
type HandlerRegistry struct {
	UploaderHandler *file.UploaderHandler
	UserHandler     *handler.UserHandler
}

func NewHandlerRegistry(
	uploaderHandler *file.UploaderHandler,
	userHandler *handler.UserHandler,
) *HandlerRegistry {
	return &HandlerRegistry{
		UploaderHandler: uploaderHandler,
		UserHandler:     userHandler,
	}
}
