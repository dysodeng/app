package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/websocket"
	webSocket "github.com/dysodeng/app/internal/interfaces/websocket"
)

// Server WebSocket服务
type Server struct {
	config   *config.Config
	ws       *webSocket.WebSocket
	wsServer *http.Server
}

// NewServer 创建WebSocket服务
func NewServer(cfg *config.Config, ws *webSocket.WebSocket) *Server {
	return &Server{
		config: cfg,
		ws:     ws,
	}
}

func (s *Server) IsEnabled() bool {
	return s.config.Server.WebSocket.Enabled
}

func (s *Server) Name() string {
	return "WebSocket"
}

// Start 启动WebSocket服务
func (s *Server) Start() error {
	// websocket 客户端连接hub
	websocket.HubBus = websocket.NewHub()
	websocket.InitMetrics()
	go websocket.HubBus.Run()

	websocket.HubBus.SetTextMessageHandler(s.ws.TextMessageHandler())
	websocket.HubBus.SetBinaryMessageHandler(s.ws.BinaryMessageHandler())

	s.wsServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.config.Server.WebSocket.Host, s.config.Server.WebSocket.Port),
		ReadHeaderTimeout: 3 * time.Second,
	}

	// 注册路由
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/v1/message", func(w http.ResponseWriter, r *http.Request) {
		websocket.Serve(w, r)
	})

	s.wsServer.Handler = mux

	var errChan = make(chan error, 1)
	go func() {
		if err := s.wsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	return fmt.Sprintf("%s:%d", s.config.Server.WebSocket.Host, s.config.Server.WebSocket.Port)
}

// Stop 停止WebSocket服务
func (s *Server) Stop(ctx context.Context) error {
	return s.wsServer.Shutdown(ctx)
}
