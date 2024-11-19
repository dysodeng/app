package router

import (
	"net/http"

	"github.com/dysodeng/app/internal/api/http/middleware"
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/service/reply/api"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// base middleware
	router.Use(middleware.CrossDomain, middleware.StartTrace)

	// api路由
	baseApiRouter := router.Group("/api/v1")

	// debug路由
	if config.App.Env != config.Prod {
		debugRouter(baseApiRouter)
	}

	// 公共组件路由
	commonRouter(baseApiRouter)

	apiRouter := baseApiRouter.Group("")
	{
		apiRouter.POST("test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, api.Success(ctx, "hello world"))
		})
	}

	return router
}
