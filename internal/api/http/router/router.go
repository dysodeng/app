package router

import (
	"net/http"

	"github.com/dysodeng/app/internal/di"

	"github.com/dysodeng/app/internal/api/http/dto/response/api"
	"github.com/dysodeng/app/internal/api/http/middleware"
	"github.com/dysodeng/app/internal/config"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// base middleware
	router.Use(middleware.CrossDomain, middleware.StartTrace)

	// api路由
	baseApiRouter := router.Group("/v1")

	appApi, err := di.InitAPI()
	if err != nil {
		panic(err)
	}

	// debug路由
	if config.App.Env != config.Prod {
		debugRouter(baseApiRouter, appApi)
	}

	// 公共组件路由
	commonRouter(baseApiRouter, appApi)

	// 文件上传路由组
	uploaderRouter := baseApiRouter.Group("/file")
	{
		uploaderRouter.POST("upload", appApi.FileUploaderController.UploadFile)
		uploaderRouter.POST("/upload/multipart/init", appApi.FileUploaderController.InitMultipartUpload)
		uploaderRouter.POST("/upload/multipart/part", appApi.FileUploaderController.UploadPart)
		uploaderRouter.POST("/upload/multipart/complete", appApi.FileUploaderController.CompleteMultipartUpload)
		uploaderRouter.POST("/upload/multipart/status", appApi.FileUploaderController.MultipartUploadStatus)
	}

	apiRouter := baseApiRouter.Group("")
	{
		apiRouter.POST("test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, api.Success(ctx, "hello world"))
		})
	}

	return router
}
