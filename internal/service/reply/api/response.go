package api

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Response api 响应数据结构
type Response[T any] struct {
	// Code 错误码
	Code Code `json:"code"`
	// Data data payload
	Data T `json:"data,omitempty"`
	// Message 错误信息
	Message string `json:"message"`
	// TraceId 追踪id
	TraceId string `json:"trace_id"`
}

// Record 分页列表记录结构
type Record[T any] struct {
	Record T     `json:"record"`
	Total  int64 `json:"total"`
}

// Success 正确响应
func Success[T any](ctx context.Context, result T) Response[T] {
	var traceId string
	if ctx.Value("Trace-Id") != nil {
		traceId = ctx.Value("Trace-Id").(string)
	} else {
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().HasTraceID() {
			traceId = span.SpanContext().TraceID().String()
		}
	}
	return Response[T]{
		Code:    CodeOk,
		Data:    result,
		Message: "success",
		TraceId: traceId,
	}
}

// Fail 失败响应
func Fail(ctx context.Context, error string, code Code) Response[any] {
	var traceId string
	if ctx.Value("Trace-Id") != nil {
		traceId = ctx.Value("Trace-Id").(string)
	} else {
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().HasTraceID() {
			traceId = span.SpanContext().TraceID().String()
		}
	}
	return Response[any]{
		Code:    code,
		Data:    nil,
		Message: error,
		TraceId: traceId,
	}
}
