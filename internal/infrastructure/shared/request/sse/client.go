package sse

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Event 表示一个解析后的SSE事件
// 兼容：标准SSE与payload内含event的风格
type Event struct {
	ID    string          // 来自 id: 行
	Event string          // 来自 event: 行；若缺失则默认为 "message"；若payload含 event 字段则覆盖
	Data  json.RawMessage // 来自 data: 多行拼接后的完整JSON字符串
	Retry int             // 来自 retry: 行（毫秒）
}

// connectOption 连接选项
// Method 默认为 GET；如需POST流式可设置 Body 与 Headers
type connectOption struct {
	Method  string
	Headers map[string]string
	Body    io.Reader
}

type ConnectOption func(*connectOption)

func WithMethod(method string) ConnectOption {
	return func(o *connectOption) {
		o.Method = method
	}
}

func WithHeaders(headers map[string]string) ConnectOption {
	return func(o *connectOption) {
		o.Headers = headers
	}
}

func WithBody(body io.Reader) ConnectOption {
	return func(o *connectOption) {
		o.Body = body
	}
}

// Client SSE 客户端
type Client struct {
	httpClient *http.Client
}

func NewSSEClient() *Client {
	return &Client{httpClient: &http.Client{}}
}

// Connect 建立SSE连接并解析事件
// 返回：事件通道、错误通道、取消函数
func (c *Client) Connect(ctx context.Context, url string, opts ...ConnectOption) (<-chan Event, <-chan error, func(), error) {
	opt := &connectOption{
		Method: http.MethodGet,
	}
	for _, o := range opts {
		o(opt)
	}

	req, err := http.NewRequest(opt.Method, url, opt.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	req = req.WithContext(ctx)

	// 必要头
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer func() {
			_ = resp.Body.Close()
		}()

		if resp.StatusCode == http.StatusBadRequest {
			b, _ := io.ReadAll(resp.Body)
			return nil, nil, nil, errors.New(string(b))
		}
		return nil, nil, nil, errors.New("unexpected status: " + resp.Status)
	}

	evCh := make(chan Event, 128)
	errCh := make(chan error, 1)

	// 取消函数：关闭body即可终止读取
	cancel := func() { _ = resp.Body.Close() }

	go func() {
		defer close(evCh)
		defer close(errCh)
		defer func() {
			_ = resp.Body.Close()
		}()

		reader := bufio.NewReader(resp.Body)
		var (
			dataLines []string
			eventName string
			id        string
			retry     int
		)

		commit := func() {
			if len(dataLines) == 0 && eventName == "" && id == "" && retry == 0 { // 空消息忽略
				return
			}
			payload := strings.Join(dataLines, "\n")
			// 标准：默认message；兼容payload内的event字段
			resolved := eventName
			if resolved == "" {
				resolved = "message"
			}
			if payload != "" {
				var probe map[string]any
				if json.Unmarshal([]byte(payload), &probe) == nil {
					if v, ok := probe["event"].(string); ok && v != "" {
						resolved = v
					}
				}
			}
			evCh <- Event{ID: id, Event: resolved, Data: json.RawMessage(payload), Retry: retry}
			// 重置状态
			dataLines = dataLines[:0]
			eventName = ""
			id = ""
			retry = 0
		}

		for {
			var line string
			line, err = reader.ReadString('\n')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					errCh <- err
				}
				return
			}
			line = strings.TrimRight(line, "\r\n")

			if line == "" { // 消息边界
				commit()
				continue
			}
			if strings.HasPrefix(line, ":") {
				// 注释心跳，忽略
				continue
			}
			// 解析字段名与值
			colon := strings.IndexByte(line, ':')
			if colon < 0 {
				continue
			}
			field := line[:colon]
			value := strings.TrimSpace(line[colon+1:])

			switch field {
			case "data":
				dataLines = append(dataLines, value)
			case "event":
				eventName = value
			case "id":
				id = value
			case "retry":
				if n, _ := strconv.Atoi(value); n > 0 {
					retry = n
				}
			default:
				// 其他字段忽略
			}
		}
	}()

	return evCh, errCh, cancel, nil
}
