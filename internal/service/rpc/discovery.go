package rpc

import (
	"strings"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/rpc"
	rpcConfig "github.com/dysodeng/rpc/config"
	"github.com/dysodeng/rpc/naming/etcd"
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var builder resolver.Builder
var discovery rpc.ServiceDiscovery

func init() {
	conf := &rpcConfig.EtcdConfig{
		Endpoints:   strings.Split(config.Etcd.Grpc.Addr, ","),
		DialTimeout: 5,
		Namespace:   config.App.Name,
	}
	builder = etcd.NewEtcdBuilder(conf, etcd.WithBuilderNamespace(config.App.Name))
	resolver.Register(builder)

	discovery = rpc.NewServiceDiscovery(builder)
}

func ServiceDiscovery() rpc.ServiceDiscovery {
	return discovery
}

func Error(err error) (error, int32) {
	grpcStatus := status.Convert(err)
	return errors.New(grpcStatus.Message()), int32(grpcStatus.Code())
}
