package di

import "github.com/google/wire"

// InfrastructureSet 基础设施
var InfrastructureSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideDB,
	ProvideRedis,
)
