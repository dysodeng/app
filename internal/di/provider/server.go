package provider

import (
	"context"

	"github.com/dysodeng/mq/contract"
	"go.uber.org/zap"

	diEvent "github.com/dysodeng/app/internal/di/event"
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	eventServer "github.com/dysodeng/app/internal/infrastructure/server/event"
	"github.com/dysodeng/app/internal/infrastructure/server/grpc"
	"github.com/dysodeng/app/internal/infrastructure/server/health"
	"github.com/dysodeng/app/internal/infrastructure/server/http"
	"github.com/dysodeng/app/internal/infrastructure/server/websocket"
	GRPC "github.com/dysodeng/app/internal/interfaces/grpc"
	HTTP "github.com/dysodeng/app/internal/interfaces/http"
	webSocket "github.com/dysodeng/app/internal/interfaces/websocket"
)

// ProvideHTTPServer 提供HTTP服务器
func ProvideHTTPServer(cfg *config.Config, handlerRegistry *HTTP.HandlerRegistry) *http.Server {
	return http.NewServer(cfg, handlerRegistry)
}

// ProvideGRPCServer 提供gRPC服务器
func ProvideGRPCServer(ctx context.Context, cfg *config.Config, serviceRegistry *GRPC.ServiceRegistry) *grpc.Server {
	return grpc.NewServer(ctx, cfg, serviceRegistry)
}

// ProvideWebSocketServer 提供WebSocket服务器
func ProvideWebSocketServer(cfg *config.Config, ws *webSocket.WebSocket) *websocket.Server {
	return websocket.NewServer(cfg, ws)
}

// ProvideHealthServer 提供容器环境健康检查服务
func ProvideHealthServer(cfg *config.Config) *health.Server {
	return health.NewServer(cfg)
}

// ProvideEventConsumerService 提供事件消费者服务
func ProvideEventConsumerService(mq contract.MQ, logger *zap.Logger) *event.ConsumerService {
	return event.NewEventConsumerService(mq.Consumer(), logger)
}

// ProvideEventServer 提供Event服务器
func ProvideEventServer(
	cfg *config.Config,
	eventConsumer *event.ConsumerService,
	registry *diEvent.HandlerRegistry,
) *eventServer.Server {
	return eventServer.NewEventServer(cfg, eventConsumer, registry)
}
