package grpc

import (
	"fmt"
	"log"
	"time"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/api/grpc/service"
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/rpc"
	"github.com/dysodeng/rpc/naming/etcd"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// grpcServer gRPC服务
type grpcServer struct {
	rpcServer rpc.Server
}

func NewGrpcServer() server.Interface {
	return &grpcServer{}
}

func (grpcServer *grpcServer) IsEnabled() bool {
	return config.Server.Grpc.Enabled
}

func (grpcServer *grpcServer) register() {
	err := grpcServer.rpcServer.RegisterService(service.NewUserService(), proto.RegisterUserServiceServer)
	if err != nil {
		log.Fatalf("grpc service register fiald: %+v\n", err)
	}
}

func (grpcServer *grpcServer) Serve() {
	log.Println("start grpc server...")
	opts := []etcd.RegistryOption{
		etcd.WithRegistryNamespace(config.App.Name),
		etcd.WithRegistryLease(10),
		etcd.WithRegistryEtcdDialTimeout(5 * time.Second),
	}
	if config.Etcd.Grpc.Username != "" {
		opts = append(opts, etcd.WithRegistryEtcdAuth(config.Etcd.Grpc.Username, config.Etcd.Grpc.Password))
	}

	registry, err := etcd.NewEtcdRegistry(
		fmt.Sprintf("%s:%s", helper.GetLocalIp(), config.Server.Grpc.Port),
		config.Etcd.Grpc.Addr,
		opts...,
	)
	if err != nil {
		log.Fatalf("grpc etcd connent fiald: %+v\n", err)
	}

	grpcServer.rpcServer = rpc.NewServer(
		config.App.Name,
		fmt.Sprintf("0.0.0.0:%s", config.Server.Grpc.Port),
		registry,
		rpc.WithServerGrpcServerOption(
			grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(trace.TracerProvider()))),
		),
	)

	// 注册服务
	grpcServer.register()

	go func() {
		if err = grpcServer.rpcServer.Serve(); err != nil {
			log.Fatalf("grpc server start fiald: %+v\n", err)
		}
	}()

	log.Printf("grpc service listening and serving 0.0.0.0:%s\n", config.Server.Grpc.Port)
}

func (grpcServer *grpcServer) Shutdown() {
	log.Println("shutdown grpc server...")
	err := grpcServer.rpcServer.Stop()
	if err != nil {
		log.Printf("grpc server shutdown fiald:%s", err)
	}
}
