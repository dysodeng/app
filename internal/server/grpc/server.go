package grpc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dysodeng/app/internal/di"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/helper"
	telemetryMetrics "github.com/dysodeng/app/internal/pkg/telemetry/metrics"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/rpc"
	rpcConfig "github.com/dysodeng/rpc/config"
	"github.com/dysodeng/rpc/logger"
	"github.com/dysodeng/rpc/metrics"
	"github.com/dysodeng/rpc/naming/etcd"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// grpcServer gRPC服务
type grpcServer struct {
	rpcServer rpc.Server
}

func NewServer() server.Server {
	return &grpcServer{}
}

func (grpcServer *grpcServer) IsEnabled() bool {
	return config.Server.Grpc.Enabled
}

func (grpcServer *grpcServer) register() {
	// 注册gRPC服务
	grpcService := di.InitGRPC()
	err := grpcServer.rpcServer.RegisterService(grpcService.UserService, proto.RegisterUserServiceServer)
	if err != nil {
		log.Fatalf("grpc service register fiald: %+v\n", err)
	}
}

func (grpcServer *grpcServer) Serve() {
	log.Println("start grpc server...")
	opts := []etcd.RegistryOption{
		etcd.WithRegistryLease(10),
		etcd.WithRegistryEtcdDialTimeout(5 * time.Second),
	}
	if config.Etcd.Grpc.Username != "" {
		opts = append(opts, etcd.WithRegistryEtcdAuth(config.Etcd.Grpc.Username, config.Etcd.Grpc.Password))
	}

	conf := &rpcConfig.ServerConfig{
		ServiceAddr: fmt.Sprintf("%s:%s", helper.GetLocalIp(), config.Server.Grpc.Port),
		EtcdConfig: rpcConfig.EtcdConfig{
			Endpoints:   strings.Split(config.Etcd.Grpc.Addr, ","),
			DialTimeout: 5,
			Namespace:   config.Server.Grpc.Namespace,
		},
	}
	if config.Monitor.Metrics.OtlpEnabled {
		conf.OtelCollectorEndpoint = config.Monitor.Metrics.OtlpEndpoint
	}

	logger.Init(config.App.Env == config.Prod)

	// 设置 meter 到 RPC 框架
	err := metrics.SetMeter(telemetryMetrics.Meter(), config.App.Name)
	if err != nil {
		log.Fatalf("grpc metrics set fiald: %+v\n", err)
	}

	registry, err := etcd.NewEtcdRegistry(conf, opts...)
	if err != nil {
		log.Fatalf("grpc etcd connent fiald: %+v\n", err)
	}

	grpcServer.rpcServer = rpc.NewServer(
		conf,
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
