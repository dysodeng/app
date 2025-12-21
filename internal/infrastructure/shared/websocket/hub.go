package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/metric"

	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/metrics"
	"github.com/dysodeng/app/internal/infrastructure/shared/websocket/message"
)

var (
	// connectionsTotal WebSocket连接总数
	connectionsTotal metric.Int64Counter
	// messagesReceivedTotal 接收消息总数
	messagesReceivedTotal metric.Int64Counter
	// messagesSentTotal 发送消息总数
	messagesSentTotal metric.Int64Counter
	// errorsTotal 错误总数
	errorsTotal metric.Int64Counter
)

// bufferSize 通道缓冲区、map 初始化大小
const bufferSize = 128

// Handler 错误处理函数
type Handler func(msg message.WsMessage, err error)

// TextMessageHandler 文本消息处理器
type TextMessageHandler interface {
	Handler(ctx context.Context, clientId, userId string, data []byte) error
	mustTextMessageHandler()
}

// BinaryMessageHandler 二进制消息处理器
type BinaryMessageHandler interface {
	Handler(ctx context.Context, clientId, userId string, data []byte) error
	mustBinaryMessageHandler()
}

type UnimplementedTextMessageHandler struct{}

func (h *UnimplementedTextMessageHandler) Handler(ctx context.Context, clientId, userId string, data []byte) error {
	panic("method Handler not implemented")
}

func (h *UnimplementedTextMessageHandler) mustTextMessageHandler() {}

type UnimplementedBinaryMessageHandler struct{}

func (h *UnimplementedBinaryMessageHandler) Handler(ctx context.Context, clientId, userId string, data []byte) error {
	panic("method Handler not implemented")
}

func (h *UnimplementedBinaryMessageHandler) mustBinaryMessageHandler() {}

var HubBus *Hub

// Hub 维护了所有的客户端连接
type Hub struct {
	// 注册请求
	register chan *Client
	// 取消注册请求
	unregister chan *Client
	// 记录 uid 跟 client 的对应关系
	userClients map[string]*Client
	// 互斥锁，保护 userClients 以及 clients 的读写
	sync.RWMutex
	// 错误处理器
	errorHandler Handler
	// 验证器
	authenticator Authenticator
	// 等待发送的消息数量
	pending atomic.Int64
	// 消息处理器
	textMessageHandler   TextMessageHandler
	binaryMessageHandler BinaryMessageHandler
}

func (h *Hub) SetTextMessageHandler(handler TextMessageHandler) {
	h.textMessageHandler = handler
}

func (h *Hub) SetBinaryMessageHandler(handler BinaryMessageHandler) {
	h.binaryMessageHandler = handler
}

// 默认的错误处理器
func defaultErrorHandler(msg message.WsMessage, err error) {
	res, _ := json.Marshal(msg)
	fmt.Printf("send message: %s, error: %s\n", string(res), err.Error())
}

func NewHub() *Hub {
	return &Hub{
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		userClients:   make(map[string]*Client, bufferSize),
		RWMutex:       sync.RWMutex{},
		errorHandler:  defaultErrorHandler,
		authenticator: &JWTAuthenticator{},
	}
}

// Run 注册、取消注册请求处理
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.userClients[client.clientId] = client
			h.Unlock()
		case client := <-h.unregister:
			h.Lock()
			client.heartbeatTicker.Stop()
			close(client.send)
			delete(h.userClients, client.clientId)
			h.Unlock()
		}
	}
}

// InitMetrics 初始化 WebSocket 监控指标
func InitMetrics() {
	meter := metrics.Meter()
	if meter == nil {
		return
	}

	// WebSocket 连接数
	_, _ = meter.Int64ObservableGauge(
		"websocket.server.connections",
		metric.WithDescription("Current number of websocket connections"),
		metric.WithUnit("{connection}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			if HubBus != nil {
				HubBus.RLock()
				count := len(HubBus.userClients)
				HubBus.RUnlock()
				o.Observe(int64(count))
			}
			return nil
		}),
	)

	// 等待发送的消息数量
	_, _ = meter.Int64ObservableGauge(
		"websocket.server.pending_messages",
		metric.WithDescription("Number of pending messages in the hub"),
		metric.WithUnit("{message}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			if HubBus != nil {
				val := HubBus.pending.Load()
				o.Observe(val)
			}
			return nil
		}),
	)

	// 计数器初始化
	connectionsTotal, _ = meter.Int64Counter(
		"websocket.server.connections_total",
		metric.WithDescription("Total number of websocket connections started"),
		metric.WithUnit("{connection}"),
	)
	messagesReceivedTotal, _ = meter.Int64Counter(
		"websocket.server.messages.received_total",
		metric.WithDescription("Total number of messages received"),
		metric.WithUnit("{message}"),
	)
	messagesSentTotal, _ = meter.Int64Counter(
		"websocket.server.messages.sent_total",
		metric.WithDescription("Total number of messages sent"),
		metric.WithUnit("{message}"),
	)
	errorsTotal, _ = meter.Int64Counter(
		"websocket.server.errors_total",
		metric.WithDescription("Total number of errors"),
		metric.WithUnit("{error}"),
	)
}

// IsClientConnected 添加客户端连接状态检查方法
func (h *Hub) IsClientConnected(clientId string) bool {
	h.RLock()
	defer h.RUnlock()
	_, exists := h.userClients[clientId]
	return exists
}
