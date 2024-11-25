//go:build wireinject
// +build wireinject

package common

import (
	commonDomain "github.com/dysodeng/app/internal/service/domain/common"
	"github.com/google/wire"
)

var AreaAppServiceSet = wire.NewSet(NewAreaAppService, commonDomain.AreaDomainServiceSet)
var ValidCodeAppServiceSet = wire.NewSet(NewValidCodeAppService, commonDomain.ValidCodeDomainServiceSet)

func InitAreaAppService() AreaAppService {
	wire.Build(AreaAppServiceSet)
	return &areaAppService{}
}

func InitValidCodeAppService() ValidCodeAppService {
	wire.Build(ValidCodeAppServiceSet)
	return &validCodeAppService{}
}
