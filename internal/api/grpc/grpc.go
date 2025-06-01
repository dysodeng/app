package grpc

import "github.com/dysodeng/app/internal/api/grpc/service"

// GRPC rpc聚合器
type GRPC struct {
	UserService *service.UserService
}

func NewGRPC(
	userService *service.UserService,
) *GRPC {
	return &GRPC{
		UserService: userService,
	}
}
