package common

import (
	"net/http"

	commonRequest "github.com/dysodeng/app/internal/api/http/request/common"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/service/app/common"
	"github.com/dysodeng/app/internal/service/reply/api"
	"github.com/gin-gonic/gin"
)

// SendValidCode 发送验证码
// @router /api/v1/common/valid_code/send [POST]
func SendValidCode(ctx *gin.Context) {
	var body commonRequest.SendValidCodeBody
	_ = ctx.ShouldBindJSON(&body)

	spanCtx := trace.Gin(ctx)

	if body.Type == "" {
		ctx.JSON(http.StatusOK, api.Fail(ctx, "缺少账号类型", api.CodeFail))
		return
	}

	var account string
	if body.Type == "sms" {
		account = body.Telephone
	} else {
		account = body.Email
	}

	validCodeAppService := common.InitValidCodeAppService()
	err := validCodeAppService.SendValidCode(spanCtx, body.Type, body.BizType, account)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, true))
}

// VerifyValidCode 验证验证码
// @router /api/v1/common/valid_code/verify [POST]
func VerifyValidCode(ctx *gin.Context) {
	var body commonRequest.VerifyValidCodeBody
	_ = ctx.ShouldBindJSON(&body)

	spanCtx := trace.Gin(ctx)

	if body.Type == "" {
		ctx.JSON(http.StatusOK, api.Fail(ctx, "缺少账号类型", api.CodeFail))
		return
	}

	var account string
	if body.Type == "sms" {
		account = body.Telephone
	} else {
		account = body.Email
	}

	validCodeAppService := common.InitValidCodeAppService()
	err := validCodeAppService.VerifyValidCode(spanCtx, body.Type, body.BizType, account, body.ValidCode)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, true))
}
