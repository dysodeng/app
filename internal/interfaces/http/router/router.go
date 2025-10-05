package router

import (
	"github.com/gin-gonic/gin"

	"github.com/dysodeng/app/internal/interfaces/http"
)

func SetupRouter(router *gin.Engine, registry *http.HandlerRegistry) {
	api := router.Group("v1")
	{
		users := api.Group("/users")
		{
			users.POST("", registry.UserHandler.Register)
			users.GET("/:id", registry.UserHandler.GetUser)
			users.GET("", registry.UserHandler.ListUsers)
			users.DELETE("/:id", registry.UserHandler.DeleteUser)
		}

		file := api.Group("file")
		{
			file.POST("/upload", registry.UploaderHandler.UploadFile)
			file.POST("/upload/multipart/init", registry.UploaderHandler.InitMultipartUpload)
			file.POST("/upload/multipart/part", registry.UploaderHandler.UploadPart)
			file.POST("/upload/multipart/complete", registry.UploaderHandler.CompleteMultipartUpload)
			file.POST("/upload/multipart/status", registry.UploaderHandler.MultipartUploadStatus)
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
