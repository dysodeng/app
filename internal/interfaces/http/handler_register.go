package http

import (
	"github.com/dysodeng/app/internal/interfaces/http/handler"
	"github.com/gin-gonic/gin"
)

// RegisterHandlers 注册所有HTTP处理器
func RegisterHandlers(engine *gin.Engine, handlers ...interface{}) {
	// 转换为Handler接口
	var httpHandlers []handler.Handler
	for _, h := range handlers {
		if httpHandler, ok := h.(handler.Handler); ok {
			httpHandlers = append(httpHandlers, httpHandler)
		}
	}
	
	// 注册所有处理器
	handler.RegisterHandlers(engine, httpHandlers...)
}