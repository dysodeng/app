package http

import (
	"github.com/gin-gonic/gin"

	"github.com/dysodeng/app/internal/interfaces/http/handler"
)

// RegisterHandlers 注册所有HTTP处理器
func RegisterHandlers(engine *gin.Engine, handlers ...handler.Handler) {
	// 转换为Handler接口
	var httpHandlers []handler.Handler
	for _, h := range handlers {
		httpHandlers = append(httpHandlers, h)
	}

	// 注册所有处理器
	handler.RegisterHandlers(engine, httpHandlers...)
}
