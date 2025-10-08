package health

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

// Server 容器环境健康检查服务
type Server struct {
	config   *config.Config
	listener net.Listener
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) IsEnabled() bool {
	return s.config.Server.Health.Enabled
}

func (s *Server) Name() string {
	return "Health"
}

func (s *Server) Addr() string {
	return fmt.Sprintf("127.0.0.1:%d", s.config.Server.Health.Port)
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp4", fmt.Sprintf(":%d", s.config.Server.Health.Port))
	if err != nil {
		return err
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				break
			}
			go s.health(conn)
		}
	}()

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	return s.listener.Close()
}

func (s *Server) health(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	_, _ = conn.Write([]byte{114, 117, 110, 110, 105, 110, 103}) // write string "running"
}
