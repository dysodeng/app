package common

import (
	"net/http"
	"net/url"

	"github.com/dysodeng/app/internal/pkg/telemetry/trace"

	"github.com/dysodeng/app/internal/pkg/qrcode"
	"github.com/dysodeng/app/internal/service/reply/api"
	"github.com/gin-gonic/gin"
)

// QrCode 生成二维码图片
// @route /api/v1/common/qr_code
func QrCode(ctx *gin.Context) {
	isUrl := ctx.DefaultQuery("is_url", "")
	text := ctx.DefaultQuery("text", "")

	spanCtx := trace.Gin(ctx)

	if text == "" {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "缺少二维码内容", api.CodeFail))
		return
	}
	if isUrl == "1" {
		text, _ = url.QueryUnescape(text)
	}

	qr, err := qrcode.NewQrCode(text, 20)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "二维码生成失败", api.CodeFail))
		return
	}

	buf, err := qr.SaveToBuffer()
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, "二维码生成失败", api.CodeFail))
		return
	}

	_, _ = ctx.Writer.Write(buf.Bytes())
}
