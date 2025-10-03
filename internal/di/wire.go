//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

// InitApp 初始化应用程序
func InitApp() (*App, error) {
	panic(wire.Build(
		ProvideConfig,
		ProvideDB,
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
