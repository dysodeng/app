package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

// Server WebSocket服务
type Server struct {
	config     *config.Config
	upgrader   websocket.Upgrader
	httpServer *http.Server
	clients    sync.Map
}

// NewServer 创建WebSocket服务
func NewServer(config *config.Config) *Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &Server{
		config:   config,
		upgrader: upgrader,
		httpServer: &http.Server{
			Addr: fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port+1), // WebSocket使用HTTP端口+1
		},
		clients: sync.Map{},
	}
}

func (s *Server) IsEnabled() bool {
	return true
}

func (s *Server) Name() string {
	return "WebSocket"
}

// Start 启动WebSocket服务
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	s.httpServer.Handler = mux

	var errChan = make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// Addr 获取服务地址
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.config.WebSocket.Host, s.config.WebSocket.Port)
}

// Stop 停止WebSocket服务
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// handleWebSocket 处理WebSocket连接
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 生成客户端ID
	clientID := r.RemoteAddr

	// 存储客户端连接
	client := &Client{
		ID:   clientID,
		Conn: conn,
	}
	s.clients.Store(clientID, client)
	defer s.clients.Delete(clientID)

	// 处理消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// 处理接收到的消息
		s.handleMessage(client, messageType, message)
	}
}

// handleMessage 处理接收到的消息
func (s *Server) handleMessage(client *Client, messageType int, message []byte) {
	// 这里可以根据业务需求处理消息
	// 示例：简单回显消息
	client.Conn.WriteMessage(messageType, message)
}

// BroadcastMessage 广播消息给所有客户端
func (s *Server) BroadcastMessage(messageType int, message []byte) {
	s.clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		client.Conn.WriteMessage(messageType, message)
		return true
	})
}

// Client WebSocket客户端
type Client struct {
	ID   string
	Conn *websocket.Conn
}
