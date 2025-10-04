package modules

import (
	"github.com/google/wire"

	fileApplicationService "github.com/dysodeng/app/internal/application/file/service"
	"github.com/dysodeng/app/internal/domain/file/service"
	fileRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/file"
	"github.com/dysodeng/app/internal/interfaces/http/handler"
	"github.com/dysodeng/app/internal/interfaces/http/handler/file"
)

// FileModule 文件模块
type FileModule struct {
	uploaderHandler *file.UploaderHandler
}

func NewFileModule(uploaderHandler *file.UploaderHandler) *FileModule {
	return &FileModule{uploaderHandler: uploaderHandler}
}

func (m *FileModule) Handlers() []handler.Handler {
	return []handler.Handler{
		m.uploaderHandler,
	}
}

func (m *FileModule) GRPCServices() []interface{} {
	return []interface{}{}
}

func (m *FileModule) EventHandlers() []interface{} {
	return []interface{}{}
}

// FileModuleWireSet 文件模块依赖注入集合
var FileModuleWireSet = wire.NewSet(
	// 仓储层
	fileRepository.NewFileRepository,
	fileRepository.NewUploaderRepository,

	// 领域层
	service.NewFileDomainService,
	service.NewUploaderDomainService,

	// 应用层
	fileApplicationService.NewUploaderApplicationService,

	// http接口层
	file.NewUploaderHandler,

	NewFileModule,
)
