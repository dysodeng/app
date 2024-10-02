package router

import (
	"net/http"

	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/trace"

	"github.com/dysodeng/app/internal/server/http/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// base middleware
	router.Use(middleware.CrossDomain, middleware.StartTrace)

	// api路由
	apiRouter := router.Group("/api/v1")
	{
		apiRouter.POST("test", func(ctx *gin.Context) {
			logger.Debug(ctx, "test")
			go func() {
				traceCtx := trace.New().NewSpan(ctx, "hello")
				logger.Info(traceCtx, "hello world")

				go func() {
					childTraceCtx := trace.New().NewSpan(traceCtx, "child")
					logger.Info(childTraceCtx, "hello child")
				}()
			}()
			ctx.JSON(http.StatusOK, api.Success(ctx, "hello world"))
		})
	}

	return router
}
