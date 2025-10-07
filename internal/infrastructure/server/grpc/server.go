package grpc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dysodeng/rpc"
	rpcConfig "github.com/dysodeng/rpc/config"
	"github.com/dysodeng/rpc/logger"
	"github.com/dysodeng/rpc/metrics"
	"github.com/dysodeng/rpc/naming/etcd"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/helper"
	telemetryMetrics "github.com/dysodeng/app/internal/infrastructure/shared/telemetry/metrics"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

// Server gRPC服务
type Server struct {
	config    *config.Config
	rpcServer rpc.Server
}

// NewServer 创建gRPC服务
func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

// Server 获取gRPC服务器
func (s *Server) Server() rpc.Server {
	return s.rpcServer
}

func (s *Server) IsEnabled() bool {
	return s.config.Server.GRPC.Enabled
}

func (s *Server) Name() string {
	return "gRPC"
}

func (s *Server) register() {}

// Start 启动gRPC服务
func (s *Server) Start() error {
	opts := []etcd.RegistryOption{
		etcd.WithRegistryLease(10),
		etcd.WithRegistryEtcdDialTimeout(5 * time.Second),
	}
	if s.config.Etcd.GRPC.Username != "" {
		opts = append(opts, etcd.WithRegistryEtcdAuth(s.config.Etcd.GRPC.Username, s.config.Etcd.GRPC.Password))
	}

	conf := &rpcConfig.ServerConfig{
		ServiceAddr: fmt.Sprintf("%s:%d", helper.GetLocalIp(), s.config.Server.GRPC.Port),
		EtcdConfig: rpcConfig.EtcdConfig{
			Endpoints:   strings.Split(s.config.Etcd.GRPC.Addr, ","),
			DialTimeout: 5,
			Namespace:   s.config.Server.GRPC.Namespace,
		},
	}
	if s.config.Monitor.Metrics.OtlpEnabled {
		conf.OtelCollectorEndpoint = s.config.Monitor.Metrics.OtlpEndpoint
	}

	logger.Init(s.config.App.Environment == config.Prod)

	// 设置 meter 到 RPC 框架
	err := metrics.SetMeter(telemetryMetrics.Meter(), s.config.App.Name)
	if err != nil {
		log.Fatalf("grpc metrics set fiald: %+v\n", err)
	}

	registry, err := etcd.NewEtcdRegistry(conf, opts...)
	if err != nil {
		log.Fatalf("grpc etcd connent fiald: %+v\n", err)
	}

	s.rpcServer = rpc.NewServer(
		conf,
		registry,
		rpc.WithServerGrpcServerOption(
			grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(trace.TracerProvider()))),
		),
	)

	// 注册服务
	s.register()

	var errChan = make(chan error, 1)
	go func() {
		if err = s.rpcServer.Serve(); err != nil {
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
	return s.rpcServer.Stop()
}
