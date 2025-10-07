package modules

import (
	"github.com/google/wire"

	"github.com/dysodeng/app/internal/application/file/event/handler"
	fileApplicationService "github.com/dysodeng/app/internal/application/file/service"
	"github.com/dysodeng/app/internal/domain/file/service"
	fileRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/file"
	"github.com/dysodeng/app/internal/interfaces/http/handler/file"
)

// FileModuleSet 文件模块依赖注入聚合
var FileModuleSet = wire.NewSet(
	// 仓储层
	fileRepository.NewFileRepository,
	fileRepository.NewUploaderRepository,

	// 领域层
	service.NewFileDomainService,
	service.NewUploaderDomainService,

	// 应用层
	fileApplicationService.NewUploaderApplicationService,

	// 事件处理层
	handler.NewFileUploadedHandler,

	// http接口层
	file.NewUploaderHandler,
)
