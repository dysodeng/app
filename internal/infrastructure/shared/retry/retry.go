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

	currentRetry := 0
	nextTry := time.Now().Add(options.waitTimeFunc(currentRetry))

	for {
		log.Printf("第%d次执行", currentRetry+1)
		err := tryFunc()
		if err == nil {
			break
		}

		currentRetry++

		if currentRetry >= options.retryNum {
			break
		}

		nextWaitTime := options.waitTimeFunc(currentRetry)
		nextTry = nextTry.Add(nextWaitTime)

		time.Sleep(time.Until(nextTry))
	}
}
