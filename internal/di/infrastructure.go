package di

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/di/event"
	"github.com/dysodeng/app/internal/di/provider"
	"github.com/dysodeng/app/internal/interfaces/grpc"
	"github.com/dysodeng/app/internal/interfaces/http"
	"github.com/dysodeng/app/internal/interfaces/websocket"
)

// InfrastructureSet 基础设施
var InfrastructureSet = wire.NewSet(
	// 基础设施
	provider.ProvideConfig,
	provider.ProvideMonitor,
	provider.ProvideLogger,
	provider.ProvideDB,
	provider.ProvideRedis,
	provider.ProvideMessageQueue,
	provider.ProvideStorage,
	provider.ProvideEventBus,
	provider.ProvideEventConsumerService,

	// 端口适配器
	provider.ProvideFileStoragePort,
	provider.ProvideFilePolicyPort,
	provider.ProvideEventPublisherPort,
	provider.ProvideTransactionManagerPort,
)

// WebSocketSet WebSocket聚合依赖
var WebSocketSet = wire.NewSet(
	websocket.NewTextMessageHandler,
	websocket.NewBinaryMessageHandler,
	websocket.NewWebSocket,
	provider.ProvideWebSocketServer,
)

// ServerSet 服务聚合依赖
var ServerSet = wire.NewSet(
	WebSocketSet,
	http.NewHandlerRegistry,
	event.NewHandlerRegistry,
	grpc.NewServiceRegistry,
	provider.ProvideHTTPServer,
	provider.ProvideGRPCServer,
	provider.ProvideHealthServer,
	provider.ProvideEventServer,
)
