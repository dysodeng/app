package user

import (
	"github.com/dysodeng/app/internal/api/grpc/proto"
	rpcService "github.com/dysodeng/app/internal/service/rpc"
	"github.com/dysodeng/rpc"
	"google.golang.org/grpc"
)

var (
	userServiceConn *grpc.ClientConn
)

// Service grpc-用户服务
func Service() (proto.UserServiceClient, error) {
	if userServiceConn == nil {
		conn, err := rpcService.ServiceDiscovery().ServiceConn(
			"user.UserService",
			rpc.WithServiceDiscoveryLB(rpc.RoundRobin),
		)
		if err != nil {
			return nil, err
		}
		userServiceConn = conn
	}
	return proto.NewUserServiceClient(userServiceConn), nil
}
