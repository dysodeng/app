package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/dysodeng/app/internal/infrastructure/shared/websocket/message"
)

// bufferSize 通道缓冲区、map 初始化大小
const bufferSize = 128

// Handler 错误处理函数
type Handler func(msg message.WsMessage, err error)

// TextMessageHandler 文本消息处理器
type TextMessageHandler interface {
	Handler(ctx context.Context, clientId, userId string, messageType int, message []byte) error
}

// BinaryMessageHandler 二进制消息处理器
type BinaryMessageHandler interface {
	Handler(ctx context.Context, clientId, userId string, data []byte) error
}

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

// Metrics 返回 Hub 的当前的关键指标
func Metrics(w http.ResponseWriter) {
	pending := HubBus.pending.Load()
	connections := len(HubBus.userClients)
	_, _ = fmt.Fprintf(w, "# HELP connections 连接数\n# TYPE connections gauge\nconnections %d\n", connections)
	_, _ = fmt.Fprintf(w, "# HELP pending 等待发送的消息数量\n# TYPE pending gauge\npending %d\n", pending)
}

// IsClientConnected 添加客户端连接状态检查方法
func (h *Hub) IsClientConnected(clientId string) bool {
	h.RLock()
	defer h.RUnlock()
	_, exists := h.userClients[clientId]
	return exists
}
