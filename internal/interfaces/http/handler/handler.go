package handler

import "github.com/gin-gonic/gin"

// Handler HTTP处理器接口
type Handler interface {
	// RegisterRoutes 注册路由
	RegisterRoutes(r *gin.RouterGroup)
}

// RegisterHandlers 注册所有处理器
func RegisterHandlers(r *gin.Engine, handlers ...Handler) {
	// API路由组
	api := r.Group("/v1")

	// 注册所有处理器
	for _, h := range handlers {
		h.RegisterRoutes(api)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
