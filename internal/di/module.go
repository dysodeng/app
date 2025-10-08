package di

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/di/modules"
)

// AllModulesSet 所有业务模块的聚合Wire Set
var AllModulesSet = wire.NewSet(
	// 在这里添加所有业务模块的Wire Set
	// 这样在wire.go中只需要引用这一个AllModulesSet
	modules.SharedModuleSet,
	modules.PassportModuleSet,
	modules.FileModuleSet,
)
