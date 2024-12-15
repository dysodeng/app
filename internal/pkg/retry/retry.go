package retry

import (
	"log"
	"time"
)

// Invoke 重试
func Invoke(tryFunc func() error, opts ...Option) {
	options := defaultRetryOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	// 重试次数
	retry := 0

	// 记录下一次尝试的时间
	nextTry := time.Now().Add(options.waitTime)

	err := tryFunc()
	if err == nil { // 请求成功
		return
	}

	for {
		log.Printf("第%d次请求", retry+1)
		// 如果当前时间大于下一次重试的时间，则等待结束，进行下一次请求
		if time.Now().After(nextTry) {
			err = tryFunc()
			if err == nil {
				break
			}
			nextWaitTime := options.waitTimeFunc(retry)
			nextTry = nextTry.Add(nextWaitTime) // 更新下一次重试的时间
		}

		if retry >= options.retryNum {
			break
		}

		retry++

		// 等待到下一次尝试的时间
		time.Sleep(nextTry.Sub(time.Now()))
	}
}
