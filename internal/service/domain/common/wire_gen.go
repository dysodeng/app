// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package common

import (
	"github.com/dysodeng/app/internal/dal/dao/common"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitAreaDomainService() AreaDomainService {
	areaDao := common.NewAreaDao()
	commonAreaDomainService := NewAreaDomainService(areaDao)
	return commonAreaDomainService
}

func InitMailDomainService() MailDomainService {
	mailDao := common.NewMailDao()
	commonMailDomainService := NewMailDomainService(mailDao)
	return commonMailDomainService
}

func InitSmsDomainService() SmsDomainService {
	smsDao := common.NewSmsDao()
	commonSmsDomainService := NewSmsDomainService(smsDao)
	return commonSmsDomainService
}

func InitValidCodeDomainService() ValidCodeDomainService {
	smsDao := common.NewSmsDao()
	commonSmsDomainService := NewSmsDomainService(smsDao)
	mailDao := common.NewMailDao()
	commonMailDomainService := NewMailDomainService(mailDao)
	commonValidCodeDomainService := NewValidCodeDomainService(commonSmsDomainService, commonMailDomainService)
	return commonValidCodeDomainService
}

// wire.go:

var AreaDomainServiceSet = wire.NewSet(common.NewAreaDao, NewAreaDomainService)

var MailDomainServiceSet = wire.NewSet(common.NewMailDao, NewMailDomainService)

var SmsDomainServiceSet = wire.NewSet(common.NewSmsDao, NewSmsDomainService)

var ValidCodeDomainServiceSet = wire.NewSet(SmsDomainServiceSet, MailDomainServiceSet, NewValidCodeDomainService)
