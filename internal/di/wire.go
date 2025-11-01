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
		InfrastructureSet, // 基础设施
		ModulesSet,        // 业务模块
		ServerSet,         // 服务模块
		NewApp,
	))
	return nil, nil
}
