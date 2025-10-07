//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"

	"github.com/dysodeng/app/internal/di/event"
	"github.com/dysodeng/app/internal/interfaces/grpc"
	"github.com/dysodeng/app/internal/interfaces/http"
)

// InitApp 初始化应用程序
func InitApp(ctx context.Context) (*App, error) {
	panic(wire.Build(
		InfrastructureSet, // 基础设施
		AllModulesSet,     // 业务模块
		WebSocketSet,      // WebSocket处理模块
		http.NewHandlerRegistry,
		event.NewHandlerRegistry,
		grpc.NewServiceRegistry,
		ProvideHTTPServer,
		ProvideGRPCServer,
		ProvideWebSocketServer,
		ProvideEventServer,
		NewApp,
	))
	return nil, nil
}
