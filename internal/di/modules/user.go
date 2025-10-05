package modules

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/application/service"
	domainService "github.com/dysodeng/app/internal/domain/service"
	"github.com/dysodeng/app/internal/infrastructure/persistence/repository"
	"github.com/dysodeng/app/internal/interfaces/http/handler"
)

// UserModuleSet 用户模块依赖注入聚合
var UserModuleSet = wire.NewSet(
	// 仓储层
	repository.NewUserRepository,

	// 领域层
	domainService.NewUserService,

	// 应用层
	service.NewUserAppService,

	// http接口层
	handler.NewUserHandler,
)
