package http

import "github.com/dysodeng/app/internal/api/http/controller/common"

// API api聚合器
type API struct {
	AreaController *common.AreaController
}

func NewAPI(
	areaController *common.AreaController,
) *API {
	return &API{
		AreaController: areaController,
	}
}
