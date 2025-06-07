package http

import (
	"github.com/dysodeng/app/internal/api/http/controller/common"
	"github.com/dysodeng/app/internal/api/http/controller/debug"
	"github.com/dysodeng/app/internal/api/http/controller/file"
	"github.com/dysodeng/app/internal/infrastructure/event/manager"
)

// API api聚合器
type API struct {
	eventManager           *manager.EventManager
	AreaController         *common.AreaController
	ValidCodeController    *common.ValidCodeController
	FileUploaderController *file.UploaderController
	DebugController        *debug.Controller
}

func NewAPI(
	eventManager *manager.EventManager,
	areaController *common.AreaController,
	validCodeController *common.ValidCodeController,
	FileUploaderController *file.UploaderController,
	debugController *debug.Controller,
) *API {
	return &API{
		eventManager:           eventManager,
		AreaController:         areaController,
		ValidCodeController:    validCodeController,
		FileUploaderController: FileUploaderController,
		DebugController:        debugController,
	}
}
