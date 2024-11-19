package middleware

import (
	"context"
	"fmt"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// StartTrace trace
func StartTrace(ctx *gin.Context) {
	savedCtx := ctx.Request.Context()
	defer func() {
		ctx.Request = ctx.Request.WithContext(savedCtx)
	}()

	traceIdByHex := ctx.Request.Header.Get("traceId")
	spanIdByHex := ctx.Request.Header.Get("spanId")

	var newCtx context.Context
	if traceIdByHex != "" {
		traceId, _ := oteltrace.TraceIDFromHex(traceIdByHex)
		spanId, _ := oteltrace.SpanIDFromHex(spanIdByHex)
		spanCtx := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID:    traceId,
			SpanID:     spanId,
			TraceFlags: oteltrace.FlagsSampled,
			Remote:     true,
		})
		carrier := propagation.HeaderCarrier{}
		carrier.Set("traceId", traceIdByHex)
		newCtx = oteltrace.ContextWithRemoteSpanContext(otel.GetTextMapPropagator().Extract(savedCtx, carrier), spanCtx)
	} else {
		newCtx = otel.GetTextMapPropagator().Extract(savedCtx, propagation.HeaderCarrier(ctx.Request.Header))
	}

	opts := []oteltrace.SpanStartOption{
		oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", ctx.Request)...),
		oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(ctx.Request)...),
		oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(config.App.Name, ctx.FullPath(), ctx.Request)...),
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
	}
	spanName := ctx.FullPath()
	if spanName == "" {
		spanName = fmt.Sprintf("HTTP %s route not found", ctx.Request.Method)
	}
	spanCtx, span := trace.Tracer().Start(newCtx, spanName, opts...)
	defer span.End()

	ctx.Request = ctx.Request.WithContext(spanCtx)

	ctx.Next()

	status := ctx.Writer.Status()
	attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
	spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(status, oteltrace.SpanKindServer)
	span.SetAttributes(attrs...)
	span.SetStatus(spanStatus, spanMessage)
	if len(ctx.Errors) > 0 {
		span.SetAttributes(attribute.String("gin.errors", ctx.Errors.String()))
	}
}
