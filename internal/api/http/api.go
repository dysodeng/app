package http

import "github.com/dysodeng/app/internal/api/http/controller/common"

// API api聚合器
type API struct {
	AreaController      *common.AreaController
	ValidCodeController *common.ValidCodeController
}

func NewAPI(
	areaController *common.AreaController,
	validCodeController *common.ValidCodeController,
) *API {
	return &API{
		AreaController:      areaController,
		ValidCodeController: validCodeController,
	}
}
