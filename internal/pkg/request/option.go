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

// WithTracer 添加链路追踪
func WithTracer(traceIdKey, spanIdKey string) Option {
	return optionFunc(func(option *requestOption) {
		option.tracerTransmit = true
		option.traceIdKey = traceIdKey
		option.traceSpanIdKey = spanIdKey
	})
}

// retryRequestOption 重试选项
type retryRequestOption struct {
	retryNum        int           // 重试次数
	initialWaitTime time.Duration // 初始等待时间
	incrementTime   time.Duration // 递增时间间隔
}

type RetryOption interface {
	apply(option *retryRequestOption)
}

type retryOptionFunc func(option *retryRequestOption)

func (f retryOptionFunc) apply(option *retryRequestOption) {
	f(option)
}

func defaultRetryOptions() *retryRequestOption {
	return &retryRequestOption{
		retryNum:        3,
		initialWaitTime: 5 * time.Second,
		incrementTime:   10 * time.Second,
	}
}

func WithRetryNum(retryNum int) RetryOption {
	return retryOptionFunc(func(option *retryRequestOption) {
		option.retryNum = retryNum
	})
}

func WithRetryInitialWaitTime(initialWaitTime time.Duration) RetryOption {
	return retryOptionFunc(func(option *retryRequestOption) {
		option.initialWaitTime = initialWaitTime
	})
}

func WithRetryIncrementTime(incrementTime time.Duration) RetryOption {
	return retryOptionFunc(func(option *retryRequestOption) {
		option.incrementTime = incrementTime
	})
}
