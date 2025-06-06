package grpc

import (
	"github.com/dysodeng/app/internal/api/grpc/service"
	"github.com/dysodeng/app/internal/infrastructure/event/manager"
)

// GRPC rpc聚合器
type GRPC struct {
	eventManager *manager.EventManager
	UserService  *service.UserService
}

func NewGRPC(
	eventManager *manager.EventManager,
	userService *service.UserService,
) *GRPC {
	return &GRPC{
		eventManager: eventManager,
		UserService:  userService,
	}
}
