//go:build wireinject
// +build wireinject

package common

import (
	commonDao "github.com/dysodeng/app/internal/dal/dao/common"
	"github.com/google/wire"
)

var AreaDomainServiceSet = wire.NewSet(commonDao.NewAreaDao, NewAreaDomainService)
var MailDomainServiceSet = wire.NewSet(commonDao.NewMailDao, NewMailDomainService)
var SmsDomainServiceSet = wire.NewSet(commonDao.NewSmsDao, NewSmsDomainService)
var ValidCodeDomainServiceSet = wire.NewSet(SmsDomainServiceSet, MailDomainServiceSet, NewValidCodeDomainService)

func InitAreaDomainService() AreaDomainService {
	wire.Build(AreaDomainServiceSet)
	return &areaDomainService{}
}

func InitMailDomainService() MailDomainService {
	wire.Build(MailDomainServiceSet)
	return &mailDomainService{}
}

func InitSmsDomainService() SmsDomainService {
	wire.Build(SmsDomainServiceSet)
	return &smsDomainService{}
}

func InitValidCodeDomainService() ValidCodeDomainService {
	wire.Build(ValidCodeDomainServiceSet)
	return &validCodeDomainService{}
}
