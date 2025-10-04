package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

// Server gRPC服务
type Server struct {
	config     *config.Config
	grpcServer *grpc.Server
}

// NewServer 创建gRPC服务
func NewServer(config *config.Config) *Server {
	grpcServer := grpc.NewServer()

	return &Server{
		config:     config,
		grpcServer: grpcServer,
	}
}

// Server 获取gRPC服务器
func (s *Server) Server() *grpc.Server {
	return s.grpcServer
}

func (s *Server) IsEnabled() bool {
	return s.config.Server.GRPC.Enabled
}

func (s *Server) Name() string {
	return "gRPC"
}

// Start 启动gRPC服务
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.GRPC.Host, s.config.Server.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	var errChan = make(chan error, 1)
	go func() {
		if err = s.grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case err = <-errChan:
		return err
	default:
		return nil
	}
}

// Addr 获取服务地址
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.config.Server.GRPC.Host, s.config.Server.GRPC.Port)
}

// Stop 停止gRPC服务
func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}
