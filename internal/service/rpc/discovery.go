package rpc

import (
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/rpc"
	"github.com/dysodeng/rpc/naming/etcd"
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var builder resolver.Builder
var discovery rpc.ServiceDiscovery

func init() {
	builder = etcd.NewEtcdBuilder(config.Etcd.Grpc.Addr, etcd.WithBuilderNamespace(config.App.Name))
	resolver.Register(builder)

	discovery = rpc.NewServiceDiscovery(config.App.Name, builder)
}

func ServiceDiscovery() rpc.ServiceDiscovery {
	return discovery
}

func Error(err error) (error, int32) {
	grpcStatus := status.Convert(err)
	return errors.New(grpcStatus.Message()), int32(grpcStatus.Code())
}
