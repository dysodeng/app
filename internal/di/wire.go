//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"
)

// InitApp 初始化应用程序
func InitApp(ctx context.Context) (*App, error) {
	panic(wire.Build(
		ProvideConfig,
		ProvideDB,
		ProvideRedis,
		AllModulesSet, // 业务模块
		ProvideModuleRegistry,
		ProvideHTTPServer,
		ProvideGRPCServer,
		ProvideWebSocketServer,
		ProvideEventBus,
		NewApp,
	))
	return nil, nil
}
