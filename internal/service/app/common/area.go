package common

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/service/domain/common"
	commonReply "github.com/dysodeng/app/internal/service/reply/common"
	"github.com/pkg/errors"
)

// AreaAppService 地区应用服务
type AreaAppService interface {
	Area(ctx context.Context, areaType string, parentAreaId string) ([]commonReply.Area, error)
	CascadeArea(ctx context.Context, provinceAreaId string, cityAreaId string, countyAreaId string) (*commonReply.CascadeArea, error)
}

// areaAppService 地区应用服务
type areaAppService struct {
	baseTraceSpanName string
	areaDomainService common.AreaDomainService
}

var areaAppServiceInstance AreaAppService

func NewAreaAppService(areaDomainService common.AreaDomainService) AreaAppService {
	if areaAppServiceInstance == nil {
		areaAppServiceInstance = &areaAppService{
			baseTraceSpanName: "service.app.common.AreaAppService",
			areaDomainService: areaDomainService,
		}
	}
	return areaAppServiceInstance
}

// Area 获取地区列表
func (area *areaAppService) Area(ctx context.Context, areaType string, parentAreaId string) ([]commonReply.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, area.baseTraceSpanName+".Area")
	defer span.End()

	if !helper.Contain(areaType, []string{"province", "city", "county"}) {
		return nil, errors.New("地区类型错误")
	}

	areaList, err := area.areaDomainService.Area(spanCtx, areaType, parentAreaId)
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
func (area *areaAppService) CascadeArea(
	ctx context.Context,
	provinceAreaId string,
	cityAreaId string,
	countyAreaId string,
) (*commonReply.CascadeArea, error) {
	spanCtx, span := trace.Tracer().Start(ctx, area.baseTraceSpanName+".CascadeArea")
	defer span.End()

	if provinceAreaId == "" {
		return nil, errors.New("缺少省地区编号")
	}
	if cityAreaId == "" {
		return nil, errors.New("缺少市地区编号")
	}
	if countyAreaId == "" {
		return nil, errors.New("缺少区县地区编号")
	}

	cascadeArea, err := area.areaDomainService.CascadeArea(spanCtx, provinceAreaId, cityAreaId, countyAreaId)
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
