package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dysodeng/app/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"

	"github.com/dysodeng/app/internal/infrastructure/config"
	ifaceHttp "github.com/dysodeng/app/internal/interfaces/http"
	"github.com/dysodeng/app/internal/interfaces/http/handler"
)

// Server HTTP服务
type Server struct {
	config     *config.Config
	engine     *gin.Engine
	httpServer *http.Server
}

// NewServer 创建HTTP服务
func NewServer(config *config.Config, handlers ...handler.Handler) *Server {
	// 设置gin模式
	if !config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())
	engine.Use(middleware.StartTrace())

	// 注册处理器
	ifaceHttp.RegisterHandlers(engine, handlers...)

	return &Server{
		config: config,
		engine: engine,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", config.Server.HTTP.Host, config.Server.HTTP.Port),
			Handler: engine,
		},
	}
}

// Engine 获取gin引擎
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) IsEnabled() bool {
	return s.config.Server.HTTP.Enabled
}

// Addr 获取服务地址
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.config.Server.HTTP.Host, s.config.Server.HTTP.Port)
}

func (s *Server) Name() string {
	return "HTTP"
}

// Start 启动HTTP服务
func (s *Server) Start() error {
	var errChan = make(chan error, 1)
	go func() {
		// ListenAndServe 只有在服务器关闭或发生错误时才会返回
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

// Stop 停止HTTP服务
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
