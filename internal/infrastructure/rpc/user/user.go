package user

import (
	"context"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	rpcService "github.com/dysodeng/app/internal/infrastructure/rpc"
	"github.com/dysodeng/rpc"
	"github.com/dysodeng/rpc/breaker"
	"github.com/dysodeng/rpc/retry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var (
	userServiceConn *grpc.ClientConn
)

// Service grpc-用户服务
func Service(ctx context.Context) (proto.UserServiceClient, error) {
	if userServiceConn == nil {
		span := trace.SpanFromContext(ctx)
		conn, err := rpcService.ServiceDiscovery().ServiceConn(
			"user.UserService",
			rpc.WithServiceDiscoveryLB(rpc.RoundRobin),
			rpc.WithServiceDiscoveryBreaker(breaker.NewCircuitBreaker()),
			rpc.WithServiceDiscoveryRetry(retry.DefaultRetryPolicy),
			rpc.WithServiceDiscoveryGrpcDialOption(
				grpc.WithStatsHandler(otelgrpc.NewClientHandler(otelgrpc.WithTracerProvider(span.TracerProvider()))),
			),
		)
		if err != nil {
			return nil, err
		}
		userServiceConn = conn
	}
	return proto.NewUserServiceClient(userServiceConn), nil
}
