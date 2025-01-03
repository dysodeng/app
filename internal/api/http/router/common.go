package router

import (
	"github.com/dysodeng/app/internal/api/http/controller/common"
	"github.com/gin-gonic/gin"
)

// commonRouter 公共组件路由
func commonRouter(router *gin.RouterGroup) {
	commonApi := router.Group("common")
	{
		commonApi.POST("area", common.Area)
		commonApi.POST("area/cascade", common.CascadeArea)
		commonApi.POST("valid_code/send", common.SendValidCode)
		commonApi.POST("valid_code/verify", common.VerifyValidCode)
		commonApi.GET("qr_code", common.QrCode)
	}
}
