package router

import (
	"github.com/gin-gonic/gin"

	"github.com/dysodeng/app/internal/interfaces/http"
)

// RegisterRouter 注册路由
func RegisterRouter(router *gin.Engine, registry *http.HandlerRegistry) {
	api := router.Group("v1")
	{
		passport := api.Group("passport")
		{
			passport.POST("login", registry.PassportHandler.Login)
			passport.POST("refresh_token", registry.PassportHandler.RefreshToken)
		}

		file := api.Group("file")
		{
			file.POST("upload", registry.UploaderHandler.UploadFile)
			file.POST("upload/multipart/init", registry.UploaderHandler.InitMultipartUpload)
			file.POST("upload/multipart/part", registry.UploaderHandler.UploadPart)
			file.POST("upload/multipart/complete", registry.UploaderHandler.CompleteMultipartUpload)
			file.POST("upload/multipart/status", registry.UploaderHandler.MultipartUploadStatus)
		}
	}

	// 健康检查
	router.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
