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
	baseApiRouter := router.Group("/api/v1")

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

	apiRouter := baseApiRouter.Group("")
	{
		apiRouter.POST("test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, api.Success(ctx, "hello world"))
		})
	}

	return router
}
