package router

import (
	"net/http"

	"github.com/dysodeng/app/internal/server/http/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// cors
	router.Use(middleware.CrossDomain)

	// api路由
	apiRouter := router.Group("/api/v1")
	{
		apiRouter.POST("test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "hello world"})
		})
	}

	return router
}
