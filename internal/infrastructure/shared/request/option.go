package request

import (
	"context"
	"time"
)

const (
	defaultRequestTimeout = 30 * time.Second
	maxBufferSize         = 512 * 1000
)

type requestOption struct {
	ctx            context.Context
	timeout        time.Duration
	maxBufferSize  int
	headers        map[string]string
	tracerTransmit bool
	traceIdKey     string
	traceSpanIdKey string
}

type Option interface {
	apply(option *requestOption)
}

type optionFunc func(option *requestOption)

func (f optionFunc) apply(option *requestOption) {
	f(option)
}

// defaultRequestOptions 默认请求选项
func defaultRequestOptions() *requestOption {
	return &requestOption{
		ctx:           context.Background(),
		timeout:       defaultRequestTimeout,
		maxBufferSize: maxBufferSize,
		headers:       make(map[string]string),
	}
}

// WithContext 设置请求上下文
func WithContext(ctx context.Context) Option {
	return optionFunc(func(option *requestOption) {
		option.ctx = ctx
	})
}

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(option *requestOption) {
		option.timeout = timeout
	})
}

// WithHeader 设置请求头
func WithHeader(key, value string) Option {
	return optionFunc(func(option *requestOption) {
		option.headers[key] = value
	})
}

// WithStreamMaxBufferSize 设置流式请求最大缓存大小
func WithStreamMaxBufferSize(maxBufferSize int) Option {
	return optionFunc(func(option *requestOption) {
		option.maxBufferSize = maxBufferSize
	})
}

// WithTraceContext 从context中自动提取trace信息
func WithTraceContext(ctx context.Context) Option {
	return optionFunc(func(option *requestOption) {
		option.ctx = ctx
		option.tracerTransmit = true
		option.traceIdKey = "X-Trace-Id"
		option.traceSpanIdKey = "X-Span-Id"
	})
}

// WithTracer 添加链路追踪
func WithTracer(traceIdKey, spanIdKey string) Option {
	return optionFunc(func(option *requestOption) {
		option.tracerTransmit = true
		option.traceIdKey = traceIdKey
		option.traceSpanIdKey = spanIdKey
	})
}
