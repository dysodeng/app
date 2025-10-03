package router

import (
	"github.com/dysodeng/app/internal/interfaces/http/handler"
	"github.com/dysodeng/app/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册HTTP路由
func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
) {
	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("", userHandler.Register)
			users.GET("/:id", userHandler.GetUser)
			users.GET("", userHandler.ListUsers)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}