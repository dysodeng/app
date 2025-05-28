//go:build wireinject
// +build wireinject

package common

import (
	commonDomain "github.com/dysodeng/app/internal/service/domain/common"
	"github.com/google/wire"
)

var ValidCodeAppServiceSet = wire.NewSet(NewValidCodeAppService, commonDomain.ValidCodeDomainServiceSet)

func InitValidCodeAppService() ValidCodeAppService {
	wire.Build(ValidCodeAppServiceSet)
	return &validCodeAppService{}
}
