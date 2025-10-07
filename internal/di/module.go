package di

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/di/modules"
	"github.com/dysodeng/app/internal/interfaces/websocket"
)

// WebSocketSet WebSocket聚合依赖
var WebSocketSet = wire.NewSet(
	websocket.NewTextMessageHandler,
	websocket.NewBinaryMessageHandler,
	websocket.NewWebSocket,
)

// AllModulesSet 所有业务模块的聚合Wire Set
var AllModulesSet = wire.NewSet(
	// 在这里添加所有业务模块的Wire Set
	// 这样在wire.go中只需要引用这一个AllModulesSet
	modules.SharedModuleSet,
	modules.PassportModuleSet,
	modules.FileModuleSet,
)
