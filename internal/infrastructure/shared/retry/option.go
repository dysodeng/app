package retry

import "time"

// option 重试选项
type option struct {
	retryNum     int                              // 重试次数
	waitTime     time.Duration                    // 重试等待时间
	waitTimeFunc func(retryNum int) time.Duration // 自定义等待时间计算函数 retryNum为当前已重试次数
}

type Option interface {
	apply(option *option)
}

type retryOptionFunc func(option *option)

func (f retryOptionFunc) apply(option *option) {
	f(option)
}

// defaultRetryOptions 默认重试选项
func defaultRetryOptions() *option {
	opt := &option{
		retryNum: 3,
		waitTime: 5 * time.Second,
	}
	opt.waitTimeFunc = func(retryNum int) time.Duration {
		return opt.waitTime
	}
	return opt
}

// WithRetryNum 设置重试次数
func WithRetryNum(retryNum int) Option {
	return retryOptionFunc(func(option *option) {
		option.retryNum = retryNum
	})
}

// WithRetryWaitTime 设置重试等待时间
func WithRetryWaitTime(waitTime time.Duration) Option {
	return retryOptionFunc(func(option *option) {
		option.waitTime = waitTime
	})
}

// WithRetryWaitTimeFunc 自定义重试等待时间计算函数
func WithRetryWaitTimeFunc(waitTimeFunc func(retryNum int) time.Duration) Option {
	return retryOptionFunc(func(option *option) {
		option.waitTimeFunc = waitTimeFunc
	})
}
