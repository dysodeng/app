package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

// Server HTTP服务
type Server struct {
	config     *config.Config
	engine     *gin.Engine
	httpServer *http.Server
}

// NewServer 创建HTTP服务
func NewServer(config *config.Config, handlers ...interface{}) *Server {
	// 设置gin模式
	if !config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	
	// 注册处理器
	RegisterHandlers(engine, handlers...)

	return &Server{
		config: config,
		engine: engine,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port),
			Handler: engine,
		},
	}
}

// Engine 获取gin引擎
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Addr 获取服务地址
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.config.HTTP.Host, s.config.HTTP.Port)
}

// Start 启动HTTP服务
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop 停止HTTP服务
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
