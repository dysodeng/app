//go:build wireinject
// +build wireinject

package user

import (
	"github.com/dysodeng/app/internal/dal/dao/user"
	"github.com/google/wire"
)

var DomainServiceSet = wire.NewSet(user.NewUserDao, NewUserDomainService)

func InitUserDomainService() DomainService {
	wire.Build(DomainServiceSet)
	return &domainService{}
}
