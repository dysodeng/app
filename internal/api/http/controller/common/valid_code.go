package common

import (
	"net/http"

	commonRequest "github.com/dysodeng/app/internal/api/http/dto/request/common"
	"github.com/dysodeng/app/internal/api/http/dto/response/api"
	"github.com/dysodeng/app/internal/application/common"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/gin-gonic/gin"
)

type ValidCodeController struct {
	baseTraceSpanName string
	validCodeService  common.ValidCodeApplicationService
}

func NewValidCodeController(validCodeService common.ValidCodeApplicationService) *ValidCodeController {
	return &ValidCodeController{
		baseTraceSpanName: "api.http.controller.common.ValidCodeController",
		validCodeService:  validCodeService,
	}
}

// SendValidCode 发送验证码
// @router /api/v1/common/valid_code/send [POST]
func (c *ValidCodeController) SendValidCode(ctx *gin.Context) {
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

	err := c.validCodeService.SendValidCode(spanCtx, body.Type, body.BizType, account)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, true))
}

// VerifyValidCode 验证验证码
// @router /api/v1/common/valid_code/verify [POST]
func (c *ValidCodeController) VerifyValidCode(ctx *gin.Context) {
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

	err := c.validCodeService.VerifyValidCode(spanCtx, body.Type, body.BizType, account, body.ValidCode)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, true))
}
