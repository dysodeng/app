package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"google.golang.org/grpc"
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

// Start 启动gRPC服务
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.GRPC.Host, s.config.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(lis)
}

// Addr 获取服务地址
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.config.GRPC.Host, s.config.GRPC.Port)
}

// Stop 停止gRPC服务
func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.GracefulStop()
}
