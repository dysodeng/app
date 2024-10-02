package api

import "context"

// Response api 响应数据结构
type Response struct {
	// Code 错误码
	Code Code `json:"code"`
	// Data data payload
	Data interface{} `json:"data,omitempty"`
	// Message 错误信息
	Message string `json:"message"`
	// TraceId 追踪id
	TraceId string `json:"trace_id"`
}

// Record 分页列表记录结构
type Record struct {
	Record          interface{} `json:"record"`
	Total           int64       `json:"total"`
	CurrentPageSize int         `json:"current_page_size"`
}

// Success 正确响应
func Success(ctx context.Context, result interface{}) Response {
	return Response{
		Code:    CodeOk,
		Data:    result,
		Message: "success",
		TraceId: ctx.Value("traceId").(string),
	}
}

// Fail 失败响应
func Fail(ctx context.Context, error string, code Code) Response {
	return Response{
		Code:    code,
		Data:    nil,
		Message: error,
		TraceId: ctx.Value("traceId").(string),
	}
}
