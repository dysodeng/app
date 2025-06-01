//go:build wireinject
// +build wireinject

package di

import (
	"github.com/dysodeng/app/internal/api/grpc"
	grpcService "github.com/dysodeng/app/internal/api/grpc/service"
	"github.com/dysodeng/app/internal/api/http"
	commonController "github.com/dysodeng/app/internal/api/http/controller/common"
	"github.com/dysodeng/app/internal/application/common"
	userAppService "github.com/dysodeng/app/internal/application/user/service"
	commonService "github.com/dysodeng/app/internal/domain/common/service"
	userService "github.com/dysodeng/app/internal/domain/user/service"
	commonRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/common"
	userRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/user"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/google/wire"
)

var (
	// 基础设施层
	InfrastructureSet = wire.NewSet(
		transactions.NewGormTransactionManager,
		commonRepository.NewAreaRepository,
		commonRepository.NewMailRepository,
		commonRepository.NewSmsRepository,
		userRepository.NewUserRepository,
	)

	// 领域层
	DomainSet = wire.NewSet(
		InfrastructureSet, // 引入基础设施依赖
		commonService.NewAreaDomainService,
		commonService.NewMailDomainService,
		commonService.NewSmsDomainService,
		commonService.NewValidCodeDomainService,
		userService.NewUserDomainService,
	)

	// 应用层
	ApplicationSet = wire.NewSet(
		DomainSet, // 引入领域层依赖
		common.NewAreaApplicationService,
		common.NewValidCodeAppService,
		userAppService.NewUserApplication,
	)

	// API聚合层
	APISet = wire.NewSet(
		ApplicationSet, // 引入应用层依赖
		commonController.NewAreaController,
		commonController.NewValidCodeController,
		http.NewAPI,
	)

	// gRPC聚合层
	GRPCSet = wire.NewSet(
		ApplicationSet,
		grpcService.NewUserService,
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
