package modules

import (
	"github.com/dysodeng/app/internal/application/passport/service"
	userDomainService "github.com/dysodeng/app/internal/domain/user/service"
	userRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/user"
	"github.com/dysodeng/app/internal/interfaces/http/handler/passport"
	"github.com/google/wire"
)

// PassportModuleSet 认证模块依赖注入聚合
var PassportModuleSet = wire.NewSet(
	// 仓储层
	userRepository.NewUserRepository,

	// 领域层
	userDomainService.NewUserDomainService,

	// 应用层
	service.NewPassportApplicationService,

	// 控制器层
	passport.NewPassportHandler,
)
