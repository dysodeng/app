package middleware

import (
	"github.com/dysodeng/app/internal/pkg/trace"
	"github.com/gin-gonic/gin"
)

// StartTrace trace
func StartTrace(ctx *gin.Context) {
	traceId := ctx.Request.Header.Get("traceId")
	parentSpanId := ctx.Request.Header.Get("parentSpanId")
	spanName := ctx.Request.Header.Get("spanName")
	spanId := trace.GenerateSpanID()

	if traceId == "" {
		traceId = spanId
	}

	ctx.Set("traceId", traceId)
	ctx.Set("spanId", spanId)
	ctx.Set("parentSpanId", parentSpanId)
	ctx.Set("parentSpanName", spanName)
	ctx.Set("spanName", ctx.FullPath())

	ctx.Next()
}
