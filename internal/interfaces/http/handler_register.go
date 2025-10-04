package http

import (
	"github.com/gin-gonic/gin"

	"github.com/dysodeng/app/internal/interfaces/http/handler"
)

// RegisterHandlers 注册所有HTTP处理器
func RegisterHandlers(engine *gin.Engine, handlers ...handler.Handler) {
	// 注册所有处理器
	handler.RegisterHandlers(engine, handlers...)
}
