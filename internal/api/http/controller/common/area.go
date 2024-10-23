package common

import (
	"net/http"

	commonRequest "github.com/dysodeng/app/internal/api/http/request/common"
	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/service/app/common"
	"github.com/gin-gonic/gin"
)

// Area 获取地区信息
// @router /api/v1/common/area [POST]
func Area(ctx *gin.Context) {
	var body commonRequest.AreaBody
	_ = ctx.ShouldBindJSON(&body)

	result, err := common.NewAreaAppService(ctx).Area(body.AreaType, body.ParentAreaId)
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

	result, err := common.NewAreaAppService(ctx).CascadeArea(body.ProvinceAreaId, body.CityAreaId, body.CountyAreaId)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, result))
}
