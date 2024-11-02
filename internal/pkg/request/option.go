package request

import (
	"context"
	"time"
)

const (
	defaultRequestTimeout = 30 * time.Second
)

type requestOption struct {
	ctx     context.Context
	timeout time.Duration
	headers map[string]string
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(option *requestOption) error

func (f optionFunc) apply(option *requestOption) error {
	return f(option)
}

func defaultRequestOptions() *requestOption {
	return &requestOption{
		ctx:     context.Background(),
		timeout: defaultRequestTimeout,
		headers: make(map[string]string),
	}
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(option *requestOption) error {
		option.ctx = ctx
		return nil
	})
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(option *requestOption) error {
		option.timeout = timeout
		return nil
	})
}

func WithHeader(key, value string) Option {
	return optionFunc(func(option *requestOption) error {
		option.headers[key] = value
		return nil
	})
}
