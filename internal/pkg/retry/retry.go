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

	// 当前重试次数
	currentRetry := 0

	// 记录下一次尝试的时间
	nextTry := time.Now().Add(options.waitTimeFunc(currentRetry))

	// err := tryFunc()
	// if err == nil { // 请求成功
	// 	return
	// }

	for {
		log.Printf("第%d次请求", currentRetry+1)
		err := tryFunc()
		if err == nil {
			break
		}

		currentRetry++

		if currentRetry >= options.retryNum {
			break
		}

		nextWaitTime := options.waitTimeFunc(currentRetry)
		nextTry = nextTry.Add(nextWaitTime) // 更新下一次重试的时间

		// 等待到下一次尝试的时间
		time.Sleep(nextTry.Sub(time.Now()))
	}
}
