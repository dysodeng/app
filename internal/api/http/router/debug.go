package router

import (
	"github.com/dysodeng/app/internal/api/http"
	"github.com/gin-gonic/gin"
)

// debugRouter debug
func debugRouter(router *gin.RouterGroup, appApi *http.API) {
	debugApi := router.Group("debug")
	{
		debugApi.POST("token", appApi.DebugController.Token)
		debugApi.POST("token/verify", appApi.DebugController.VerifyToken)
		debugApi.POST("random_string", appApi.DebugController.GenRandomString)
		debugApi.POST("gorm_logger", appApi.DebugController.GormLogger)
		debugApi.POST("grpc/user", appApi.DebugController.User)
		debugApi.POST("grpc/user/list", appApi.DebugController.ListUser)
		debugApi.POST("grpc/user/create", appApi.DebugController.CreateUser)
		debugApi.POST("chat/message", appApi.DebugController.ChatMessage)
		debugApi.POST("remote_request", appApi.DebugController.RemoteRequest)
		debugApi.POST("retry", appApi.DebugController.Retry)
		debugApi.POST("cache", appApi.DebugController.Cache)
	}
}
