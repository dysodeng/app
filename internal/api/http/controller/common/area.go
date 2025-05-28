package common

import (
	"net/http"

	commonRequest "github.com/dysodeng/app/internal/api/http/dto/request/common"
	"github.com/dysodeng/app/internal/api/http/dto/response/api"
	"github.com/dysodeng/app/internal/application/common"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/gin-gonic/gin"
)

type AreaController struct {
	baseTraceSpanName string
	areaService       common.AreaApplicationService
}

func NewAreaController(areaService common.AreaApplicationService) *AreaController {
	return &AreaController{
		baseTraceSpanName: "api.http.controller.common.AreaController",
		areaService:       areaService,
	}
}

// Area 获取地区信息
// @router /api/v1/common/area [POST]
func (c *AreaController) Area(ctx *gin.Context) {
	var body commonRequest.AreaBody
	_ = ctx.ShouldBindJSON(&body)

	spanCtx := trace.Gin(ctx)
	result, err := c.areaService.Area(spanCtx, body.AreaType, body.ParentAreaId)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, result))
}

// CascadeArea 级联获取地区信息
// @router /api/v1/common/area/cascade [POST]
func (c *AreaController) CascadeArea(ctx *gin.Context) {
	var body commonRequest.CascadeAreaBody
	_ = ctx.ShouldBindJSON(&body)

	spanCtx := trace.Gin(ctx)
	result, err := c.areaService.CascadeArea(spanCtx, body.ProvinceAreaId, body.CityAreaId, body.CountyAreaId)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, result))
}
