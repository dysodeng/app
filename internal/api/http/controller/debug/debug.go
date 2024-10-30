package debug

import (
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/api/token"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/trace"
	"github.com/gin-gonic/gin"
)

func Token(ctx *gin.Context) {
	t, _ := token.GenerateToken("user", map[string]interface{}{
		"user_id": 1,
	}, nil)
	ctx.JSON(200, api.Success(ctx, t))
}

func GormLogger(ctx *gin.Context) {
	traceCtx := trace.New().NewSpan(ctx, "debug.GormLogger")

	var mailConfig common.MailConfig
	db.DB().WithContext(traceCtx).First(&mailConfig)
	var smsConfig common.SmsConfig
	db.DB().WithContext(traceCtx).Where("a=?", "b").First(&smsConfig)

	go func() {
		childTraceCtx := trace.New().NewSpan(traceCtx, "debug.GormLogger.child")
		logger.Debug(childTraceCtx, "child logger")
		logger.Error(childTraceCtx, "child logger")
	}()

	ctx.JSON(200, api.Success(traceCtx, mailConfig))
}
