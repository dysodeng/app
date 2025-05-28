package service

import (
	"context"

	"github.com/dysodeng/app/internal/domain/common/model"
	"github.com/dysodeng/app/internal/domain/common/repository"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
)

// AreaDomainService 地区领域服务
type AreaDomainService interface {
	ProvinceByAreaId(ctx context.Context, areaId string) (*model.Area, error)
	ProvinceByAreaName(ctx context.Context, areaName string) (*model.Area, error)
	CityByAreaId(ctx context.Context, areaId string) (*model.Area, error)
	CityByAreaName(ctx context.Context, areaName string) (*model.Area, error)
	CountyByAreaId(ctx context.Context, areaId string) (*model.Area, error)
	CountyByAreaName(ctx context.Context, areaName string) (*model.Area, error)
	CascadeArea(ctx context.Context, provinceAreaId string, cityAreaId string, countyAreaId string) (*model.CascadeArea, error)
	Area(ctx context.Context, areaType string, parentAreaId string) ([]model.Area, error)
}

// areaDomainService Area领域服务
type areaDomainService struct {
	baseTraceSpanName string
	areaRepo          repository.AreaRepository
}

func NewAreaDomainService(areaRepo repository.AreaRepository) AreaDomainService {
	return &areaDomainService{
		baseTraceSpanName: "domain.common.service.AreaDomainService",
		areaRepo:          areaRepo,
	}
}

// ProvinceByAreaId 根据省地区编号查询省地区信息
// @param areaId 省地区编号
func (svc *areaDomainService) ProvinceByAreaId(ctx context.Context, areaId string) (*model.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".ProvinceByAreaId")
	defer span.End()

	area, err := svc.areaRepo.ProvinceByAreaId(spanCtx, areaId)
	if err != nil {
		logger.Error(spanCtx, "根据地区ID查询[地区-省]数据失败", logger.ErrorField(err))
		return nil, errors.New("省地区查询失败")
	}

	return model.AreaFromProvince(*area), nil
}

// ProvinceByAreaName 根据省地区名称查询省地区信息
// @param areaName 省地区名称
func (svc *areaDomainService) ProvinceByAreaName(ctx context.Context, areaName string) (*model.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".ProvinceByAreaName")
	defer span.End()

	area, err := svc.areaRepo.ProvinceByAreaName(spanCtx, areaName)
	if err != nil {
		logger.Error(spanCtx, "根据地区名称查询[地区-省]数据失败", logger.ErrorField(err))
		return nil, errors.New("省地区查询失败")
	}

	return model.AreaFromProvince(*area), nil
}

// CityByAreaId 根据市地区编号查询市地区信息
// @param areaId 市地区编号
func (svc *areaDomainService) CityByAreaId(ctx context.Context, areaId string) (*model.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CityByAreaId")
	defer span.End()

	area, err := svc.areaRepo.CityByAreaId(spanCtx, areaId, "")
	if err != nil {
		logger.Error(spanCtx, "根据地区ID查询[地区-市]数据失败", logger.ErrorField(err))
		return nil, errors.New("市地区查询失败")
	}

	return model.AreaFromCity(*area), nil
}

// CityByAreaName 根据市地区名称查询市地区信息
// @param areaName 市地区名称
func (svc *areaDomainService) CityByAreaName(ctx context.Context, areaName string) (*model.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CityByAreaName")
	defer span.End()

	area, err := svc.areaRepo.CityByAreaName(spanCtx, areaName)
	if err != nil {
		logger.Error(spanCtx, "根据地区名称查询[地区-市]数据失败", logger.ErrorField(err))
		return nil, errors.New("市地区查询失败")
	}

	return model.AreaFromCity(*area), nil
}

// CountyByAreaId 根据县地区编号查询县地区信息
// @param areaId 县地区编号
func (svc *areaDomainService) CountyByAreaId(ctx context.Context, areaId string) (*model.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CountyByAreaId")
	defer span.End()

	area, err := svc.areaRepo.CountyByAreaId(spanCtx, areaId, "")
	if err != nil {
		logger.Error(spanCtx, "根据地区ID查询[地区-区县]数据失败", logger.ErrorField(err))
		return nil, errors.New("区县地区查询失败")
	}

	return model.AreaFromCounty(*area), nil
}

// CountyByAreaName 根据区县地区名称查询区县地区信息
// @param areaName 区县地区名称
func (svc *areaDomainService) CountyByAreaName(ctx context.Context, areaName string) (*model.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CountyByAreaName")
	defer span.End()

	area, err := svc.areaRepo.CountyByAreaName(spanCtx, areaName)
	if err != nil {
		logger.Error(spanCtx, "根据地区名称查询[地区-区县]数据失败", logger.ErrorField(err))
		return nil, errors.New("区县地区查询失败")
	}

	return model.AreaFromCounty(*area), nil
}

// CascadeArea 级联查询
// @param provinceAreaId 省地区编号
// @param cityAreaId 市地区编号
// @param countyAreaId 区县地区编号
func (svc *areaDomainService) CascadeArea(ctx context.Context, provinceAreaId string, cityAreaId string, countyAreaId string) (*model.CascadeArea, error) {
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

	province, err := svc.areaRepo.ProvinceByAreaId(spanCtx, provinceAreaId)
	if err != nil {
		logger.Error(spanCtx, "查询省列表失败", logger.ErrorField(err))
		return nil, errors.New("查询省列表失败")
	}
	if province == nil || province.ID <= 0 {
		return nil, errors.New("省地区错误")
	}
	city, err := svc.areaRepo.CityByAreaId(spanCtx, cityAreaId, province.AreaId)
	if err != nil {
		logger.Error(spanCtx, "查询市列表失败", logger.ErrorField(err))
		return nil, errors.New("查询市列表失败")
	}
	if city == nil || city.ID <= 0 {
		return nil, errors.New("市地区错误")
	}
	county, err := svc.areaRepo.CountyByAreaId(spanCtx, countyAreaId, city.AreaId)
	if err != nil {
		logger.Error(spanCtx, "查询区县列表失败", logger.ErrorField(err))
		return nil, errors.New("查询区县列表失败")
	}
	if county == nil || county.ID <= 0 {
		return nil, errors.New("区县地区错误")
	}

	return model.CascadeAreaFromAreas(*province, *city, *county), nil
}

// Area 查询下级地区列表
// @param areaType 地区类型 省：province 市：city 区县：county
// @param parentAreaId 父级地区ID
func (svc *areaDomainService) Area(ctx context.Context, areaType string, parentAreaId string) ([]model.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".Area")
	defer span.End()

	if areaType == "" {
		return nil, errors.New("缺少地区类型")
	}
	if areaType != "province" && parentAreaId == "" {
		return nil, errors.New("缺少地区父ID")
	}

	var areaCollection []model.Area
	switch areaType {
	case "province":
		provinceList, err := svc.areaRepo.ProvinceList(spanCtx)
		if err != nil {
			logger.Error(spanCtx, "查询省列表失败", logger.ErrorField(err))
			return nil, errors.New("查询省列表失败")
		}
		areaCollection = model.AreaListFromProvince(provinceList)
	case "city":
		cityList, err := svc.areaRepo.CityList(spanCtx, parentAreaId)
		if err != nil {
			logger.Error(spanCtx, "查询市列表失败", logger.ErrorField(err))
			return nil, errors.New("查询市列表失败")
		}
		areaCollection = model.AreaListFromCity(cityList)
	case "county":
		countyList, err := svc.areaRepo.CountyList(spanCtx, parentAreaId)
		if err != nil {
			logger.Error(spanCtx, "查询区县列表失败", logger.ErrorField(err))
			return nil, errors.New("查询区县列表失败")
		}
		areaCollection = model.AreaListFromCounty(countyList)
	default:
		return nil, errors.New("地区类型错误")
	}

	return areaCollection, nil
}
