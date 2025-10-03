package di

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/di/modules"
	"github.com/dysodeng/app/internal/interfaces/http/handler"
)

// Module 依赖注入模块接口
type Module interface {
	// Handlers 返回模块的HTTP处理器
	Handlers() []handler.Handler
	// GRPCServices 返回模块的gRPC服务
	GRPCServices() []interface{}
	// EventHandlers 返回模块的事件处理器
	EventHandlers() []interface{}
}

// ModuleRegistrar 模块注册器
type ModuleRegistrar struct {
	modules []Module
}

// NewModuleRegistrar 创建模块注册器
func NewModuleRegistrar() *ModuleRegistrar {
	return &ModuleRegistrar{
		modules: []Module{},
	}
}

// Register 注册模块
func (r *ModuleRegistrar) Register(module Module) {
	r.modules = append(r.modules, module)
}

// GetAllModules 获取所有模块
func (r *ModuleRegistrar) GetAllModules() []Module {
	return r.modules
}

// AllModulesSet 所有业务模块的聚合Wire Set
var AllModulesSet = wire.NewSet(
	// 在这里添加所有业务模块的Wire Set
	// 这样在wire.go中只需要引用这一个AllModulesSet
	modules.UserModuleWireSet,
)

// ProvideModuleRegistry 提供模块注册表
func ProvideModuleRegistry(
	userModule *modules.UserModule,
) *ModuleRegistrar {
	registrar := NewModuleRegistrar()
	registrar.Register(userModule)
	return registrar
}
