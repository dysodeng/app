package common

import (
	"net/http"

	commonRequest "github.com/dysodeng/app/internal/api/http/request/common"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/service/app/common"
	"github.com/dysodeng/app/internal/service/reply/api"
	"github.com/gin-gonic/gin"
)

// Area 获取地区信息
// @router /api/v1/common/area [POST]
func Area(ctx *gin.Context) {
	var body commonRequest.AreaBody
	_ = ctx.ShouldBindJSON(&body)

	spanCtx := trace.Gin(ctx)

	areaAppService := common.InitAreaAppService()
	result, err := areaAppService.Area(spanCtx, body.AreaType, body.ParentAreaId)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, result))
}

// CascadeArea 级联获取地区信息
// @router /api/v1/common/area/cascade [POST]
func CascadeArea(ctx *gin.Context) {
	var body commonRequest.CascadeAreaBody
	_ = ctx.ShouldBindJSON(&body)

	spanCtx := trace.Gin(ctx)

	areaAppService := common.InitAreaAppService()
	result, err := areaAppService.CascadeArea(spanCtx, body.ProvinceAreaId, body.CityAreaId, body.CountyAreaId)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, result))
}
