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

	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
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
				span.SetAttributes(attribute.String("request.HttpRequest.error", err.Error()))
			} else {
				span.SetStatus(codes.Ok, "")
			}
			if response != nil && response.Request != nil {
				span.SetAttributes(
					attribute.Int("http.status_code", response.StatusCode),
					attribute.String("http.method", response.Request.Method),
					attribute.String("http.url", response.Request.URL.String()),
				)
			}
		}()
	}

	// 检查传入的context是否有较短的超时时间
	var timeoutCtx context.Context
	var cancel context.CancelFunc

	if deadline, ok := reqOpts.ctx.Deadline(); ok {
		// 如果传入的context有deadline，检查是否比我们设置的超时更短
		remainingTime := time.Until(deadline)
		if remainingTime < reqOpts.timeout {
			// 如果传入context的剩余时间比设置的超时更短，使用独立的context
			// 但保留trace信息
			if span := oteltrace.SpanFromContext(reqOpts.ctx); span.SpanContext().IsValid() {
				// 创建带有trace信息的新context
				newCtx := oteltrace.ContextWithSpan(context.Background(), span)
				timeoutCtx, cancel = context.WithTimeout(newCtx, reqOpts.timeout)
			} else {
				timeoutCtx, cancel = context.WithTimeout(context.Background(), reqOpts.timeout)
			}
		} else {
			timeoutCtx, cancel = context.WithTimeout(reqOpts.ctx, reqOpts.timeout)
		}
	} else {
		timeoutCtx, cancel = context.WithTimeout(reqOpts.ctx, reqOpts.timeout)
	}
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, method, requestUrl, body)
	if err != nil {
		return nil, 0, err
	}
	for headerName, headerValue := range reqOpts.headers {
		req.Header.Add(headerName, headerValue)
	}

	client := &http.Client{}
	response, err = client.Do(req)
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
		logger.Error(reqOpts.ctx, "请求错误", logger.AddField("response", string(d)))
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
			if response != nil && response.Request != nil {
				span.SetAttributes(
					attribute.Int("http.status_code", response.StatusCode),
					attribute.String("http.method", response.Request.Method),
					attribute.String("http.url", response.Request.URL.String()),
				)
			}
		}()
	}

	// 检查传入的context是否有较短的超时时间
	var timeoutCtx context.Context
	var cancel context.CancelFunc

	if deadline, ok := reqOpts.ctx.Deadline(); ok {
		// 如果传入的context有deadline，检查是否比我们设置的超时更短
		remainingTime := time.Until(deadline)
		if remainingTime < reqOpts.timeout {
			// 如果传入context的剩余时间比设置的超时更短，使用独立的context
			// 但保留trace信息
			if span := oteltrace.SpanFromContext(reqOpts.ctx); span.SpanContext().IsValid() {
				// 创建带有trace信息的新context
				newCtx := oteltrace.ContextWithSpan(context.Background(), span)
				timeoutCtx, cancel = context.WithTimeout(newCtx, reqOpts.timeout)
			} else {
				timeoutCtx, cancel = context.WithTimeout(context.Background(), reqOpts.timeout)
			}
		} else {
			timeoutCtx, cancel = context.WithTimeout(reqOpts.ctx, reqOpts.timeout)
		}
	} else {
		timeoutCtx, cancel = context.WithTimeout(reqOpts.ctx, reqOpts.timeout)
	}
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, method, requestUrl, body)
	if err != nil {
		return 0, err
	}
	for headerName, headerValue := range reqOpts.headers {
		req.Header.Add(headerName, headerValue)
	}

	client := &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 && response.StatusCode != 201 {
		d, _ := io.ReadAll(response.Body)
		logger.Error(reqOpts.ctx, "请求错误", logger.AddField("response", string(d)))
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
			logger.Error(reqOpts.ctx, "请求错误", logger.AddField("response", string(bts)))
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
	for key, val := range data {
		body.Set(key, val)
	}
	reader := bytes.NewReader([]byte(body.Encode()))
	opts = append(opts, WithHeader("Content-Type", "application/x-www-form-urlencoded"))
	return request(requestUrl, method, reader, opts...)
}
