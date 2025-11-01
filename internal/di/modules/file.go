package modules

import (
	"github.com/google/wire"

	fileDecorator "github.com/dysodeng/app/internal/application/file/decorator"
	"github.com/dysodeng/app/internal/application/file/event/handler"
	fileApplicationService "github.com/dysodeng/app/internal/application/file/service"
	fileRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/file"
	fileGRPCService "github.com/dysodeng/app/internal/interfaces/grpc/service"
	"github.com/dysodeng/app/internal/interfaces/http/handler/file"
)

// FileModuleSet 文件模块依赖注入聚合
var FileModuleSet = wire.NewSet(
	// 仓储层
	fileRepository.NewFileRepository,
	fileRepository.NewUploaderRepository,

	// 领域层
	fileDecorator.NewFileDomainServiceWithTracing,
	fileDecorator.NewUploaderDomainServiceWithTracing,

	// 应用层
	fileApplicationService.NewFileApplicationService,
	fileApplicationService.NewUploaderApplicationService,

	// 事件处理层
	handler.NewFileUploadedHandler,

	// grpc接口层
	fileGRPCService.NewFileService,

	// http接口层
	file.NewUploaderHandler,
)
