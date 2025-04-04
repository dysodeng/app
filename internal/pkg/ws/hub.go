package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/dysodeng/app/internal/pkg/ws/message"
)

// bufferSize 通道缓冲区、map 初始化大小
const bufferSize = 128

// Handler 错误处理函数
type Handler func(msg message.WsMessage, err error)

// MessageHandler 消息处理器
type MessageHandler interface {
	Handler(clientId, userId, userType string, messageType int, message []byte) error
}

var HubBus *Hub

// Hub 维护了所有的客户端连接
type Hub struct {
	// 注册请求
	register chan *Client
	// 取消注册请求
	unregister chan *Client
	// 记录 clientId 跟 client 的对应关系
	clients map[string]*Client
	// 互斥锁，保护 userClients 以及 clients 的读写
	sync.RWMutex
	// 错误处理器
	errorHandler Handler
	// 验证器
	authenticator Authenticator
	// 等待发送的消息数量
	pending atomic.Int64
	// 消息处理器
	messageHandler MessageHandler
}

func (h *Hub) SetMessageHandler(handler MessageHandler) {
	h.messageHandler = handler
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
		clients:       make(map[string]*Client, bufferSize),
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
			h.clients[client.clientId] = client
			h.Unlock()
		case client := <-h.unregister:
			h.Lock()
			client.heartbeatTicker.Stop()
			close(client.send)
			delete(h.clients, client.clientId)
			h.Unlock()
		}
	}
}

// Metrics 返回 Hub 的当前的关键指标
func Metrics(w http.ResponseWriter) {
	pending := HubBus.pending.Load()
	connections := len(HubBus.clients)
	_, _ = w.Write([]byte(fmt.Sprintf("# HELP connections 连接数\n# TYPE connections gauge\nconnections %d\n", connections)))
	_, _ = w.Write([]byte(fmt.Sprintf("# HELP pending 等待发送的消息数量\n# TYPE pending gauge\npending %d\n", pending)))
}
