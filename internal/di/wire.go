//go:build wireinject
// +build wireinject

package di

import (
	"github.com/dysodeng/app/internal/api/grpc"
	grpcService "github.com/dysodeng/app/internal/api/grpc/service"
	"github.com/dysodeng/app/internal/api/http"
	commonController "github.com/dysodeng/app/internal/api/http/controller/common"
	debugController "github.com/dysodeng/app/internal/api/http/controller/debug"
	fileController "github.com/dysodeng/app/internal/api/http/controller/file"
	"github.com/dysodeng/app/internal/application/common"
	fileEventHandler "github.com/dysodeng/app/internal/application/file/event/handler"
	fileAppService "github.com/dysodeng/app/internal/application/file/service"
	userEventHandler "github.com/dysodeng/app/internal/application/user/event/handler"
	userAppService "github.com/dysodeng/app/internal/application/user/service"
	commonService "github.com/dysodeng/app/internal/domain/common/service"
	fileDomainService "github.com/dysodeng/app/internal/domain/file/service"
	userService "github.com/dysodeng/app/internal/domain/user/service"
	"github.com/dysodeng/app/internal/infrastructure/event/bus"
	"github.com/dysodeng/app/internal/infrastructure/event/manager"
	"github.com/dysodeng/app/internal/infrastructure/event/publisher"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
	commonRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/common"
	fileRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/file"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/google/wire"
)

var (
	// 数据持久化层
	PersistenceSet = wire.NewSet(
		// 数据持久化
		transactions.NewGormTransactionManager,
		commonRepository.NewAreaRepository,
		commonRepository.NewMailRepository,
		commonRepository.NewSmsRepository,
		fileRepository.NewFileRepository,
		fileRepository.NewUploaderRepository,

		// 缓存基础设施
		cache.NewCacheFactory,
		ProvideTypedCache,

		// 用户仓储提供者
		ProvideUserRepository,
	)

	// 基础设施层
	InfrastructureSet = wire.NewSet(
		PersistenceSet, // 引入数据持久化层

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
		fileDomainService.NewFileDomainService,
		fileDomainService.NewUploaderDomainService,
	)

	// 应用层
	ApplicationSet = wire.NewSet(
		DomainSet, // 引入领域层依赖

		// 应用服务
		common.NewAreaApplicationService,
		common.NewValidCodeAppService,
		userAppService.NewUserApplicationService,
		fileAppService.NewUploaderApplicationService,

		// 事件处理器
		NewEventManagerWithHandlers,
		userEventHandler.NewUserCreatedHandler,
		fileEventHandler.NewFileUploadedHandler,
	)

	// API聚合层
	APISet = wire.NewSet(
		ApplicationSet, // 引入应用层依赖

		// Api Controller
		commonController.NewAreaController,
		commonController.NewValidCodeController,
		debugController.NewDebugController,
		fileController.NewUploaderController,

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

func InitAPI() (*http.API, error) {
	wire.Build(APISet)
	return &http.API{}, nil
}

func InitGRPC() (*grpc.GRPC, error) {
	wire.Build(GRPCSet)
	return &grpc.GRPC{}, nil
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

// ProvideTypedCache 缓存工厂
func ProvideTypedCache(factory *cache.Factory) contract.TypedCache {
	return factory.GetTypedCache()
}
