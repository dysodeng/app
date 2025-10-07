package grpc

import (
	"github.com/dysodeng/rpc"

	v1 "github.com/dysodeng/app/api/generated/go/proto/file/v1"
	"github.com/dysodeng/app/internal/infrastructure/shared/errors"
	"github.com/dysodeng/app/internal/interfaces/grpc/service"
)

// ServiceRegistry grpc服务注册表
type ServiceRegistry struct {
	fileService *service.FileService
}

func NewServiceRegistry(fileService *service.FileService) *ServiceRegistry {
	return &ServiceRegistry{
		fileService: fileService,
	}
}

func (registry *ServiceRegistry) RegisterGRPCService(srv rpc.Server) error {
	return errors.NewPipeline().Then(func() error {
		return srv.RegisterService(registry.fileService, v1.RegisterFileServiceServer)
	}).Execute()
}
