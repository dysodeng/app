package di

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/interfaces/websocket"
)

// InfrastructureSet 基础设施
var InfrastructureSet = wire.NewSet(
	ProvideConfig,
	ProvideMonitor,
	ProvideLogger,
	ProvideDB,
	ProvideRedis,
	ProvideMessageQueue,
	ProvideStorage,
	ProvideEventBus,
	ProvideEventConsumerService,
)

// WebSocketSet WebSocket聚合依赖
var WebSocketSet = wire.NewSet(
	websocket.NewTextMessageHandler,
	websocket.NewBinaryMessageHandler,
	websocket.NewWebSocket,
)
