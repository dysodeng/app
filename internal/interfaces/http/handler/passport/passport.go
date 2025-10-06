package passport

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dysodeng/app/internal/application/passport/dto/command"
	"github.com/dysodeng/app/internal/application/passport/service"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
	"github.com/dysodeng/app/internal/interfaces/http/dto/request/passport"
	"github.com/dysodeng/app/internal/interfaces/http/dto/response/api"
	"github.com/dysodeng/app/internal/interfaces/http/validator"
)

// Handler 认证
type Handler struct {
	baseTraceSpanName string
	passportService   service.PassportApplicationService
}

func NewPassportHandler(passportService service.PassportApplicationService) *Handler {
	return &Handler{
		baseTraceSpanName: "interfaces.http.handler.passport.Handler",
		passportService:   passportService,
	}
}

func (h *Handler) Login(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), h.baseTraceSpanName+".Login")
	defer span.End()

	platformType := ctx.GetHeader("PlatformType")

	var req passport.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, validator.TransError(err), api.CodeFail))
		return
	}

	res, err := h.passportService.Login(spanCtx, &command.LoginCommand{
		PlatformType: platformType,
		UserType:     req.UserType,
		GrantType:    req.GrantType,
		WxCode:       req.WxCode,
		Code:         req.Code,
		OpenId:       req.OpenId,
		Username:     req.Username,
		Password:     req.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(spanCtx, res))
}

func (h *Handler) RefreshToken(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), h.baseTraceSpanName+".RefreshToken")
	defer span.End()

	var req passport.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, validator.TransError(err), api.CodeFail))
		return
	}

	res, err := h.passportService.RefreshToken(spanCtx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(spanCtx, res))
}
