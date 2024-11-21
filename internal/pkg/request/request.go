package request

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

func request(requestUrl, method string, body io.Reader, opts ...Option) ([]byte, int, error) {
	reqOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt.apply(reqOpts)
	}

	timeoutCtx, cancel := context.WithTimeout(reqOpts.ctx, reqOpts.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, method, requestUrl, body)
	if err != nil {
		return nil, 0, err
	}
	for headerName, headerValue := range reqOpts.headers {
		req.Header.Add(headerName, headerValue)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, 0, errors.New("请求超时")
		}
		return nil, 0, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 && response.StatusCode != 201 {
		d, _ := io.ReadAll(response.Body)
		return d, response.StatusCode, errors.New("请求错误")
	}

	b, _ := io.ReadAll(response.Body)

	return b, response.StatusCode, nil
}

func streamRequest(requestUrl, method string, body io.Reader, fn func([]byte) error, opts ...Option) (int, error) {
	reqOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt.apply(reqOpts)
	}

	timeoutCtx, cancel := context.WithTimeout(reqOpts.ctx, reqOpts.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, method, requestUrl, body)
	if err != nil {
		return 0, err
	}
	for headerName, headerValue := range reqOpts.headers {
		req.Header.Add(headerName, headerValue)
	}

	client := &http.Client{Timeout: reqOpts.timeout}
	response, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 && response.StatusCode != 201 {
		d, _ := io.ReadAll(response.Body)
		return response.StatusCode, errors.New(string(d))
	}

	scanner := bufio.NewScanner(response.Body)
	// increase the buffer size to avoid running out of space
	scanBuf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(scanBuf, maxBufferSize)
	for scanner.Scan() {
		bts := scanner.Bytes()
		if err = fn(bts); err != nil {
			return 0, err
		}
	}

	select {
	case <-timeoutCtx.Done():
		return 0, context.DeadlineExceeded
	default:

	}

	return response.StatusCode, nil
}

// Request http请求
func Request(requestUrl, method string, body io.Reader, opts ...Option) ([]byte, int, error) {
	return request(requestUrl, method, body, opts...)
}

// StreamRequest 流式请求
func StreamRequest(requestUrl, method string, body io.Reader, fn func([]byte) error, opts ...Option) (int, error) {
	return streamRequest(requestUrl, method, body, fn, opts...)
}

// JsonRequest json请求
func JsonRequest(requestUrl, method string, data map[string]interface{}, opts ...Option) ([]byte, int, error) {
	dataBytes, _ := json.Marshal(data)
	opts = append(opts, WithHeader("Content-Type", "application/json; charset=utf-8"))
	return request(requestUrl, method, bytes.NewReader(dataBytes), opts...)
}

// FormRequest form-data请求
func FormRequest(requestUrl, method string, data map[string]string, opts ...Option) ([]byte, int, error) {
	body := url.Values{}
	if data != nil {
		for key, val := range data {
			body.Set(key, val)
		}
	}
	reader := bytes.NewReader([]byte(body.Encode()))
	opts = append(opts, WithHeader("Content-Type", "application/x-www-form-urlencoded"))
	return request(requestUrl, method, reader, opts...)
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
