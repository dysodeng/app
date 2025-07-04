// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/dysodeng/app/internal/api/grpc"
	service6 "github.com/dysodeng/app/internal/api/grpc/service"
	"github.com/dysodeng/app/internal/api/http"
	common3 "github.com/dysodeng/app/internal/api/http/controller/common"
	"github.com/dysodeng/app/internal/api/http/controller/debug"
	file2 "github.com/dysodeng/app/internal/api/http/controller/file"
	common2 "github.com/dysodeng/app/internal/application/common"
	handler2 "github.com/dysodeng/app/internal/application/file/event/handler"
	service5 "github.com/dysodeng/app/internal/application/file/service"
	"github.com/dysodeng/app/internal/application/user/event/handler"
	service2 "github.com/dysodeng/app/internal/application/user/service"
	service3 "github.com/dysodeng/app/internal/domain/common/service"
	service4 "github.com/dysodeng/app/internal/domain/file/service"
	"github.com/dysodeng/app/internal/domain/user/service"
	"github.com/dysodeng/app/internal/infrastructure/event/bus"
	"github.com/dysodeng/app/internal/infrastructure/event/manager"
	"github.com/dysodeng/app/internal/infrastructure/event/publisher"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
	"github.com/dysodeng/app/internal/infrastructure/persistence/repository/common"
	"github.com/dysodeng/app/internal/infrastructure/persistence/repository/file"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitAPI() (*http.API, error) {
	inMemoryEventBus := bus.NewInMemoryEventBus()
	transactionManager := transactions.NewGormTransactionManager()
	factory, err := cache.NewCacheFactory()
	if err != nil {
		return nil, err
	}
	typedCache := ProvideTypedCache(factory)
	userRepository := ProvideUserRepository(transactionManager, typedCache)
	domainEventPublisher := publisher.NewDomainEventPublisher(inMemoryEventBus)
	userDomainService := service.NewUserDomainService(userRepository, domainEventPublisher)
	userApplicationService := service2.NewUserApplicationService(userDomainService)
	userCreatedHandler := handler.NewUserCreatedHandler(userApplicationService)
	eventManager := NewEventManagerWithHandlers(inMemoryEventBus, userCreatedHandler)
	areaRepository := common.NewAreaRepository(transactionManager)
	areaDomainService := service3.NewAreaDomainService(areaRepository)
	areaApplicationService := common2.NewAreaApplicationService(areaDomainService)
	areaController := common3.NewAreaController(areaApplicationService)
	smsRepository := common.NewSmsRepository(transactionManager)
	smsDomainService := service3.NewSmsDomainService(smsRepository)
	mailRepository := common.NewMailRepository(transactionManager)
	mailDomainService := service3.NewMailDomainService(mailRepository)
	validCodeDomainService := service3.NewValidCodeDomainService(smsDomainService, mailDomainService)
	validCodeApplicationService := common2.NewValidCodeAppService(validCodeDomainService)
	validCodeController := common3.NewValidCodeController(validCodeApplicationService)
	fileRepository := file.NewFileRepository(transactionManager)
	uploaderRepository := file.NewUploaderRepository(transactionManager)
	uploaderDomainService := service4.NewUploaderDomainService(transactionManager, fileRepository, uploaderRepository, domainEventPublisher)
	uploaderApplicationService := service5.NewUploaderApplicationService(uploaderDomainService)
	uploaderController := file2.NewUploaderController(uploaderApplicationService)
	controller := debug.NewDebugController()
	api := http.NewAPI(eventManager, areaController, validCodeController, uploaderController, controller)
	return api, nil
}

func InitGRPC() (*grpc.GRPC, error) {
	inMemoryEventBus := bus.NewInMemoryEventBus()
	transactionManager := transactions.NewGormTransactionManager()
	factory, err := cache.NewCacheFactory()
	if err != nil {
		return nil, err
	}
	typedCache := ProvideTypedCache(factory)
	userRepository := ProvideUserRepository(transactionManager, typedCache)
	domainEventPublisher := publisher.NewDomainEventPublisher(inMemoryEventBus)
	userDomainService := service.NewUserDomainService(userRepository, domainEventPublisher)
	userApplicationService := service2.NewUserApplicationService(userDomainService)
	userCreatedHandler := handler.NewUserCreatedHandler(userApplicationService)
	eventManager := NewEventManagerWithHandlers(inMemoryEventBus, userCreatedHandler)
	userService := service6.NewUserService(userApplicationService)
	grpcGRPC := grpc.NewGRPC(eventManager, userService)
	return grpcGRPC, nil
}

// wire.go:

var (
	// 数据持久化层
	PersistenceSet = wire.NewSet(transactions.NewGormTransactionManager, common.NewAreaRepository, common.NewMailRepository, common.NewSmsRepository, file.NewFileRepository, file.NewUploaderRepository, cache.NewCacheFactory, ProvideTypedCache,

		ProvideUserRepository,
	)

	// 基础设施层
	InfrastructureSet = wire.NewSet(
		PersistenceSet, bus.NewInMemoryEventBus, publisher.NewDomainEventPublisher, wire.Bind(new(bus.EventBus), new(*bus.InMemoryEventBus)),
	)

	// 领域层
	DomainSet = wire.NewSet(
		InfrastructureSet, service3.NewAreaDomainService, service3.NewMailDomainService, service3.NewSmsDomainService, service3.NewValidCodeDomainService, service.NewUserDomainService, service4.NewFileDomainService, service4.NewUploaderDomainService,
	)

	// 应用层
	ApplicationSet = wire.NewSet(
		DomainSet, common2.NewAreaApplicationService, common2.NewValidCodeAppService, service2.NewUserApplicationService, service5.NewUploaderApplicationService, NewEventManagerWithHandlers, handler.NewUserCreatedHandler, handler2.NewFileUploadedHandler,
	)

	// API聚合层
	APISet = wire.NewSet(
		ApplicationSet, common3.NewAreaController, common3.NewValidCodeController, debug.NewDebugController, file2.NewUploaderController, http.NewAPI,
	)

	// gRPC聚合层
	GRPCSet = wire.NewSet(
		ApplicationSet, service6.NewUserService, grpc.NewGRPC,
	)
)

// NewEventManagerWithHandlers 注册事件处理器
func NewEventManagerWithHandlers(
	eventBus bus.EventBus,
	userCreatedHandler *handler.UserCreatedHandler,
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
