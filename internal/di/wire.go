//go:build wireinject
// +build wireinject

package di

import (
	"github.com/dysodeng/app/internal/api/http"
	commonController "github.com/dysodeng/app/internal/api/http/controller/common"
	"github.com/dysodeng/app/internal/application/common"
	"github.com/dysodeng/app/internal/domain/common/service"
	commonRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/common"
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
	)

	// 领域层
	DomainSet = wire.NewSet(
		InfrastructureSet, // 引入基础设施依赖
		service.NewAreaDomainService,
		service.NewMailDomainService,
		service.NewSmsDomainService,
		service.NewValidCodeDomainService,
	)

	// 应用层
	ApplicationSet = wire.NewSet(
		DomainSet, // 引入领域层依赖
		common.NewAreaApplicationService,
		common.NewValidCodeAppService,
	)

	// API聚合层
	APISet = wire.NewSet(
		ApplicationSet, // 引入应用层依赖
		commonController.NewAreaController,
		commonController.NewValidCodeController,
		http.NewAPI,
	)
)

func InitAPI() *http.API {
	wire.Build(APISet)
	return &http.API{}
}
