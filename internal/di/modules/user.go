package modules

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/application/service"
	domainService "github.com/dysodeng/app/internal/domain/service"
	"github.com/dysodeng/app/internal/infrastructure/persistence/repository"
	"github.com/dysodeng/app/internal/interfaces/http/handler"
)

// UserModule 用户模块
type UserModule struct {
	userHandler handler.Handler
}

// NewUserModule 创建用户模块
func NewUserModule(userHandler *handler.UserHandler) *UserModule {
	return &UserModule{
		userHandler: userHandler,
	}
}

// Handlers 返回模块的HTTP处理器
func (m *UserModule) Handlers() []handler.Handler {
	return []handler.Handler{
		m.userHandler,
	}
}

// GRPCServices 返回模块的gRPC服务
func (m *UserModule) GRPCServices() []interface{} {
	return []interface{}{}
}

// EventHandlers 返回模块的事件处理器
func (m *UserModule) EventHandlers() []interface{} {
	return []interface{}{}
}

// UserModuleWireSet 用户模块依赖注入集合
var UserModuleWireSet = wire.NewSet(
	// 仓储层
	repository.NewUserRepository,

	// 领域层
	domainService.NewUserService,

	// 应用层
	service.NewUserAppService,

	// http接口层
	handler.NewUserHandler,

	NewUserModule,
)
