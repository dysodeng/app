package grpc

import (
	"log"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
)

type grpcServer struct {
}

func NewGrpcServer() server.Interface {
	return &grpcServer{}
}

func (grpcServer *grpcServer) Serve() {
	if !config.Server.Grpc.Enabled {
		return
	}
	log.Println("start grpc server...")
}

func (grpcServer *grpcServer) Shutdown() {
	if !config.Server.Grpc.Enabled {
		return
	}
	log.Println("shutdown grpc server...")
}
