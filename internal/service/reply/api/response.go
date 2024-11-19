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
	span := trace.SpanFromContext(ctx)
	return Response[T]{
		Code:    CodeOk,
		Data:    result,
		Message: "success",
		TraceId: span.SpanContext().TraceID().String(),
	}
}

// Fail 失败响应
func Fail(ctx context.Context, error string, code Code) Response[any] {
	span := trace.SpanFromContext(ctx)
	return Response[any]{
		Code:    code,
		Data:    nil,
		Message: error,
		TraceId: span.SpanContext().TraceID().String(),
	}
}
