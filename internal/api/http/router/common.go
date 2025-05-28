package router

import (
	"github.com/dysodeng/app/internal/api/http"
	"github.com/dysodeng/app/internal/api/http/controller/common"
	"github.com/gin-gonic/gin"
)

// commonRouter 公共组件路由
func commonRouter(router *gin.RouterGroup, appApi *http.API) {
	commonApi := router.Group("common")
	{
		commonApi.POST("area", appApi.AreaController.Area)
		commonApi.POST("area/cascade", appApi.AreaController.CascadeArea)
		commonApi.POST("valid_code/send", common.SendValidCode)
		commonApi.POST("valid_code/verify", common.VerifyValidCode)
		commonApi.GET("qr_code", common.QrCode)
	}
}
