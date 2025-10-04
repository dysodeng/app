package di

import "github.com/google/wire"

// InfrastructureSet 基础设施
var InfrastructureSet = wire.NewSet(
	ProvideConfig,
	ProvideMonitor,
	ProvideLogger,
	ProvideDB,
	ProvideRedis,
	ProvideStorage,
)
