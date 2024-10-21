package router

import (
	"github.com/dysodeng/app/internal/api/http/controller/debug"
	"github.com/gin-gonic/gin"
)

// debugRouter debug
func debugRouter(router *gin.RouterGroup) {
	debugApi := router.Group("debug")
	{
		debugApi.POST("token", debug.Token)
		debugApi.POST("gorm_logger", debug.GormLogger)
	}
}
