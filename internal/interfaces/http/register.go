package http

import (
	"github.com/dysodeng/app/internal/interfaces/http/handler/file"
	"github.com/dysodeng/app/internal/interfaces/http/handler/passport"
)

// HandlerRegistry 控制器注册表
type HandlerRegistry struct {
	PassportHandler *passport.Handler
	UploaderHandler *file.UploaderHandler
}

func NewHandlerRegistry(
	passportHandler *passport.Handler,
	uploaderHandler *file.UploaderHandler,
) *HandlerRegistry {
	return &HandlerRegistry{
		PassportHandler: passportHandler,
		UploaderHandler: uploaderHandler,
	}
}
