package common

import (
	"context"

	"github.com/dysodeng/app/internal/api/http/dto/response/common"
	"github.com/dysodeng/app/internal/domain/common/service"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
)

type AreaApplicationService interface {
	Area(ctx context.Context, areaType string, parentAreaId string) ([]common.Area, error)
	CascadeArea(ctx context.Context, provinceAreaId string, cityAreaId string, countyAreaId string) (*common.CascadeArea, error)
}

type areaApplicationService struct {
	baseTraceSpanName string
	areaDomainService service.AreaDomainService
}

func NewAreaApplicationService(areaDomainService service.AreaDomainService) AreaApplicationService {
	return &areaApplicationService{
		baseTraceSpanName: "application.common.AreaApplicationService",
		areaDomainService: areaDomainService,
	}
}

func (svc *areaApplicationService) Area(ctx context.Context, areaType string, parentAreaId string) ([]common.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".Area")
	defer span.End()

	if !helper.Contain([]string{"province", "city", "county"}, areaType) {
		return nil, errors.New("地区类型错误")
	}

	areaList, err := svc.areaDomainService.Area(spanCtx, areaType, parentAreaId)
	if err != nil {
		return nil, err
	}

	var areaCollection []common.Area
	for _, areaItem := range areaList {
		areaCollection = append(areaCollection, common.Area{
			AreaId:       areaItem.AreaId,
			AreaName:     areaItem.AreaName,
			ParentAreaId: areaItem.ParentAreaId,
		})
	}

	return areaCollection, nil
}

func (svc *areaApplicationService) CascadeArea(ctx context.Context, provinceAreaId string, cityAreaId string, countyAreaId string) (*common.CascadeArea, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CascadeArea")
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

	cascadeArea, err := svc.areaDomainService.CascadeArea(spanCtx, provinceAreaId, cityAreaId, countyAreaId)
	if err != nil {
		return nil, err
	}

	return &common.CascadeArea{
		Province: common.Area{
			AreaId:   cascadeArea.Province.AreaId,
			AreaName: cascadeArea.Province.AreaName,
		},
		City: common.Area{
			AreaId:       cascadeArea.City.AreaId,
			AreaName:     cascadeArea.City.AreaName,
			ParentAreaId: cascadeArea.City.ParentId,
		},
		County: common.Area{
			AreaId:       cascadeArea.County.AreaId,
			AreaName:     cascadeArea.County.AreaName,
			ParentAreaId: cascadeArea.County.ParentId,
		},
	}, nil
}
