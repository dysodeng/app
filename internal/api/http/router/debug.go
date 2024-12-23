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
		debugApi.POST("random_string", debug.GenRandomString)
		debugApi.POST("gorm_logger", debug.GormLogger)
		debugApi.POST("grpc/user", debug.User)
		debugApi.POST("grpc/user/list", debug.ListUser)
		debugApi.POST("grpc/user/create", debug.CreateUser)
		debugApi.POST("chat/message", debug.ChatMessage)
		debugApi.POST("remote_request", debug.RemoteRequest)
		debugApi.POST("retry", debug.Retry)
		debugApi.POST("cache", debug.Cache)
	}
}
