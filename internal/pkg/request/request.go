package request

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	netUrl "net/url"
	"time"

	"github.com/dysodeng/app/internal/pkg/logger"

	"github.com/pkg/errors"
)

// Request http请求
func Request(ctx context.Context, url, method string, body io.Reader, opts ...Option) ([]byte, int, error) {
	reqOpts := defaultRequestOptions()
	for _, opt := range opts {
		_ = opt.apply(reqOpts)
	}

	timeoutCtx, cancel := context.WithTimeout(reqOpts.ctx, reqOpts.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, method, url, body)
	if err != nil {
		logger.Error(ctx, "请求错误", logger.ErrorField(err))
		return nil, 0, err
	}
	for headerName, headerValue := range reqOpts.headers {
		req.Header.Add(headerName, headerValue)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, 0, errors.New("请求超时")
		}
		logger.Error(ctx, "请求错误", logger.ErrorField(err))
		return nil, 0, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		d, _ := io.ReadAll(res.Body)
		logger.Error(
			ctx,
			"请求错误",
			logger.Field{Key: "http_status", Value: res.StatusCode},
			logger.Field{Key: "http_body", Value: string(d)},
		)
		return d, res.StatusCode, errors.New("请求错误")
	}

	b, _ := io.ReadAll(res.Body)

	return b, res.StatusCode, nil
}

// JsonRequest json请求
func JsonRequest(ctx context.Context, url, method string, data map[string]interface{}, opts ...Option) ([]byte, int, error) {
	dataBytes, _ := json.Marshal(data)
	opts = append(opts, WithHeader("Content-Type", "application/json; charset=utf-8"))
	return Request(ctx, url, method, bytes.NewReader(dataBytes), opts...)
}

// FormRequest form-data请求
func FormRequest(ctx context.Context, url, method string, data map[string]string, opts ...Option) ([]byte, int, error) {
	body := netUrl.Values{}
	if data != nil {
		for key, val := range data {
			body.Set(key, val)
		}
	}
	reader := bytes.NewReader([]byte(body.Encode()))
	opts = append(opts, WithHeader("Content-Type", "application/x-www-form-urlencoded"))
	return Request(ctx, url, method, reader, opts...)
}

// RetryRequest 带重试的请求
func RetryRequest(requestFunc func() error, retryNum int) {
	// 设置初始等待时间和递增时间间隔
	initialWait := 5 * time.Second
	increment := 10 * time.Second
	retry := 0 // 重试次数

	// 记录下一次尝试的时间
	nextTry := time.Now().Add(initialWait)

	err := requestFunc()
	if err == nil { // 请求成功
		return
	}

	for {
		// 如果当前时间大于下一次尝试的时间，则等待结束，进行下一次请求
		if time.Now().After(nextTry) {
			err = requestFunc()
			if err == nil {
				break
			}
			nextTry = nextTry.Add(increment) // 更新下一次尝试的时间
		}

		if retry >= retryNum {
			break
		}

		retry++

		// 等待到下一次尝试的时间
		time.Sleep(nextTry.Sub(time.Now()))
	}
}
