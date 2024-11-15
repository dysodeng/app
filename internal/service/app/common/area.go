package common

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/trace"
	"github.com/dysodeng/app/internal/service/domain/common"
	commonReply "github.com/dysodeng/app/internal/service/reply/common"
	"github.com/pkg/errors"
)

// AreaAppService 地区应用服务
type AreaAppService struct {
	ctx               context.Context
	areaDomainService *common.AreaDomainService
	baseTraceSpanName string
}

func NewAreaAppService(ctx context.Context) *AreaAppService {
	baseTraceSpanName := "app.service.common.AreaAppService"
	traceCtx := trace.New().NewSpan(ctx, baseTraceSpanName)
	return &AreaAppService{
		ctx:               traceCtx,
		areaDomainService: common.NewAreaDomainService(traceCtx),
		baseTraceSpanName: baseTraceSpanName,
	}
}

// Area 获取地区列表
func (area *AreaAppService) Area(areaType string, parentAreaId string) ([]commonReply.Area, error) {
	if !helper.Contain(areaType, []string{"province", "city", "county"}) {
		return nil, errors.New("地区类型错误")
	}

	areaList, err := area.areaDomainService.Area(areaType, parentAreaId)
	if err != nil {
		return nil, err
	}

	var areaCollection []commonReply.Area
	for _, areaItem := range areaList {
		areaCollection = append(areaCollection, commonReply.Area{
			AreaId:       areaItem.AreaId,
			AreaName:     areaItem.AreaName,
			ParentAreaId: areaItem.ParentAreaId,
		})
	}

	return areaCollection, nil
}

// CascadeArea 获取地区级联信息
func (area *AreaAppService) CascadeArea(
	provinceAreaId string,
	cityAreaId string,
	countyAreaId string,
) (*commonReply.CascadeArea, error) {
	if provinceAreaId == "" {
		return nil, errors.New("缺少省地区编号")
	}
	if cityAreaId == "" {
		return nil, errors.New("缺少市地区编号")
	}
	if countyAreaId == "" {
		return nil, errors.New("缺少区县地区编号")
	}

	cascadeArea, err := area.areaDomainService.CascadeArea(provinceAreaId, cityAreaId, countyAreaId)
	if err != nil {
		return nil, err
	}

	return &commonReply.CascadeArea{
		Province: commonReply.Area{
			AreaId:   cascadeArea.Province.AreaId,
			AreaName: cascadeArea.Province.AreaName,
		},
		City: commonReply.Area{
			AreaId:       cascadeArea.City.AreaId,
			AreaName:     cascadeArea.City.AreaName,
			ParentAreaId: cascadeArea.City.ParentAreaId,
		},
		County: commonReply.Area{
			AreaId:       cascadeArea.County.AreaId,
			AreaName:     cascadeArea.County.AreaName,
			ParentAreaId: cascadeArea.County.ParentAreaId,
		},
	}, nil
}
