//go:build wireinject
// +build wireinject

package di

import (
	"github.com/dysodeng/app/internal/api/grpc"
	grpcService "github.com/dysodeng/app/internal/api/grpc/service"
	"github.com/dysodeng/app/internal/api/http"
	commonController "github.com/dysodeng/app/internal/api/http/controller/common"
	debugController "github.com/dysodeng/app/internal/api/http/controller/debug"
	"github.com/dysodeng/app/internal/application/common"
	userEventHandler "github.com/dysodeng/app/internal/application/user/event/handler"
	userAppService "github.com/dysodeng/app/internal/application/user/service"
	commonService "github.com/dysodeng/app/internal/domain/common/service"
	userService "github.com/dysodeng/app/internal/domain/user/service"
	"github.com/dysodeng/app/internal/infrastructure/event/bus"
	"github.com/dysodeng/app/internal/infrastructure/event/manager"
	"github.com/dysodeng/app/internal/infrastructure/event/publisher"
	commonRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/common"
	userRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/user"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/google/wire"
)

var (
	// 基础设施层
	InfrastructureSet = wire.NewSet(
		// 数据持久化
		transactions.NewGormTransactionManager,
		commonRepository.NewAreaRepository,
		commonRepository.NewMailRepository,
		commonRepository.NewSmsRepository,
		userRepository.NewUserRepository,

		// 事件处理器
		bus.NewInMemoryEventBus,
		publisher.NewDomainEventPublisher,
		wire.Bind(new(bus.EventBus), new(*bus.InMemoryEventBus)),
	)

	// 领域层
	DomainSet = wire.NewSet(
		InfrastructureSet, // 引入基础设施依赖

		// 领域服务
		commonService.NewAreaDomainService,
		commonService.NewMailDomainService,
		commonService.NewSmsDomainService,
		commonService.NewValidCodeDomainService,
		userService.NewUserDomainService,
	)

	// 应用层
	ApplicationSet = wire.NewSet(
		DomainSet, // 引入领域层依赖

		// 应用服务
		common.NewAreaApplicationService,
		common.NewValidCodeAppService,
		userAppService.NewUserApplicationService,

		// 事件处理器
		NewEventManagerWithHandlers,
		userEventHandler.NewUserCreatedHandler,
	)

	// API聚合层
	APISet = wire.NewSet(
		ApplicationSet, // 引入应用层依赖

		// Api Controller
		commonController.NewAreaController,
		commonController.NewValidCodeController,
		debugController.NewDebugController,

		// api聚合器
		http.NewAPI,
	)

	// gRPC聚合层
	GRPCSet = wire.NewSet(
		ApplicationSet, // 引入应用层依赖

		// gRPC服务
		grpcService.NewUserService,

		// gRPC聚合器
		grpc.NewGRPC,
	)
)

func InitAPI() *http.API {
	wire.Build(APISet)
	return &http.API{}
}

func InitGRPC() *grpc.GRPC {
	wire.Build(GRPCSet)
	return &grpc.GRPC{}
}

// NewEventManagerWithHandlers 注册事件处理器
func NewEventManagerWithHandlers(
	eventBus bus.EventBus,
	userCreatedHandler *userEventHandler.UserCreatedHandler,
) *manager.EventManager {
	return manager.NewEventManager(
		eventBus,
		userCreatedHandler,
	)
}
