package request

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func request(requestUrl, method string, body io.Reader, opts ...Option) ([]byte, int, error) {
	reqOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt.apply(reqOpts)
	}

	var response *http.Response
	var err error

	var traceId, traceSpanId string
	if reqOpts.tracerTransmit {
		_, span := trace.Tracer().Start(reqOpts.ctx, "request.HttpRequest")
		defer span.End()
		if span.SpanContext().HasTraceID() {
			traceId = span.SpanContext().TraceID().String()
		}
		if span.SpanContext().HasSpanID() {
			traceSpanId = span.SpanContext().SpanID().String()
		}

		reqOpts.headers[reqOpts.traceIdKey] = traceId
		reqOpts.headers[reqOpts.traceSpanIdKey] = traceSpanId

		defer func() {
			if err != nil {
				span.SetStatus(codes.Ok, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}
			span.SetAttributes(
				attribute.Int("http.status_code", response.StatusCode),
				attribute.String("http.method", response.Request.Method),
				attribute.String("http.url", response.Request.URL.String()),
			)
		}()
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

	response, err = http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = errors.New("请求超时")
			return nil, 0, err
		}
		return nil, 0, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 && response.StatusCode != 201 {
		d, _ := io.ReadAll(response.Body)
		err = errors.New("请求错误")
		return d, response.StatusCode, err
	}

	b, _ := io.ReadAll(response.Body)

	return b, response.StatusCode, nil
}

func streamRequest(requestUrl, method string, body io.Reader, fn func([]byte) error, opts ...Option) (int, error) {
	reqOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt.apply(reqOpts)
	}

	var response *http.Response
	var err error

	var traceId, traceSpanId string
	if reqOpts.tracerTransmit {
		_, span := trace.Tracer().Start(reqOpts.ctx, "request.StreamRequest")
		defer span.End()
		if span.SpanContext().HasTraceID() {
			traceId = span.SpanContext().TraceID().String()
		}
		if span.SpanContext().HasSpanID() {
			traceSpanId = span.SpanContext().SpanID().String()
		}

		reqOpts.headers[reqOpts.traceIdKey] = traceId
		reqOpts.headers[reqOpts.traceSpanIdKey] = traceSpanId

		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}
			span.SetAttributes(
				attribute.Int("http.status_code", response.StatusCode),
				attribute.String("http.method", response.Request.Method),
				attribute.String("http.url", response.Request.URL.String()),
			)
		}()
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
	response, err = client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 && response.StatusCode != 201 {
		d, _ := io.ReadAll(response.Body)
		err = errors.New(string(d))
		return response.StatusCode, err
	}

	scanner := bufio.NewScanner(response.Body)
	// increase the buffer size to avoid running out of space
	scanBuf := make([]byte, 0, reqOpts.maxBufferSize)
	scanner.Buffer(scanBuf, reqOpts.maxBufferSize)
	for scanner.Scan() {
		bts := scanner.Bytes()
		if err = fn(bts); err != nil {
			return 0, err
		}
	}

	select {
	case <-timeoutCtx.Done():
		err = context.DeadlineExceeded
		return 0, err
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
func RetryRequest(requestFunc func() error, opts ...RetryOption) {
	options := defaultRetryOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	// 重试次数
	retry := 0

	// 记录下一次尝试的时间
	nextTry := time.Now().Add(options.initialWaitTime)

	err := requestFunc()
	if err == nil { // 请求成功
		return
	}

	for {
		log.Printf("第%d次请求", retry+1)
		// 如果当前时间大于下一次重试的时间，则等待结束，进行下一次请求
		if time.Now().After(nextTry) {
			err = requestFunc()
			if err == nil {
				break
			}
			nextTry = nextTry.Add(options.incrementTime) // 更新下一次重试的时间
		}

		if retry >= options.retryNum {
			break
		}

		retry++

		// 等待到下一次尝试的时间
		time.Sleep(nextTry.Sub(time.Now()))
	}
}
