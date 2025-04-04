package common

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/telemetry/trace"

	commonDao "github.com/dysodeng/app/internal/dal/dao/common"
	"github.com/dysodeng/app/internal/pkg/logger"
	commonDo "github.com/dysodeng/app/internal/service/do/common"
	"github.com/pkg/errors"
)

// AreaDomainService 地区领域服务
type AreaDomainService interface {
	ProvinceByAreaId(ctx context.Context, areaId string) (*commonDo.Area, error)
	ProvinceByAreaName(ctx context.Context, areaName string) (*commonDo.Area, error)
	CityByAreaId(ctx context.Context, areaId string) (*commonDo.Area, error)
	CityByAreaName(ctx context.Context, areaName string) (*commonDo.Area, error)
	CountyByAreaId(ctx context.Context, areaId string) (*commonDo.Area, error)
	CountyByAreaName(ctx context.Context, areaName string) (*commonDo.Area, error)
	CascadeArea(ctx context.Context, provinceAreaId string, cityAreaId string, countyAreaId string) (*commonDo.CascadeArea, error)
	Area(ctx context.Context, areaType string, parentAreaId string) ([]commonDo.Area, error)
}

// areaDomainService Area领域服务
type areaDomainService struct {
	baseTraceSpanName string
	areaDao           commonDao.AreaDao
}

var areaDomainServiceInstance AreaDomainService

func NewAreaDomainService(areaDao commonDao.AreaDao) AreaDomainService {
	if areaDomainServiceInstance == nil {
		areaDomainServiceInstance = &areaDomainService{
			baseTraceSpanName: "service.domain.common.AreaDomainService",
			areaDao:           areaDao,
		}
	}
	return areaDomainServiceInstance
}

// ProvinceByAreaId 根据省地区编号查询省地区信息
// @param areaId 省地区编号
func (ads *areaDomainService) ProvinceByAreaId(ctx context.Context, areaId string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".ProvinceByAreaId")
	defer span.End()

	area, err := ads.areaDao.ProvinceByAreaId(spanCtx, areaId)
	if err != nil {
		return nil, err
	}
	return &commonDo.Area{
		AreaId:   area.AreaId,
		AreaName: area.AreaName,
	}, nil
}

// CityByAreaId 根据市地区编号查询市地区信息
// @param areaId 市地区编号
func (ads *areaDomainService) CityByAreaId(ctx context.Context, areaId string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".CityByAreaId")
	defer span.End()

	area, err := ads.areaDao.CityByAreaId(spanCtx, areaId, "")
	if err != nil {
		return nil, err
	}
	return &commonDo.Area{
		AreaId:       area.AreaId,
		AreaName:     area.AreaName,
		ParentAreaId: area.ParentId,
	}, nil
}

// CountyByAreaId 根据县地区编号查询县地区信息
// @param areaId 县地区编号
func (ads *areaDomainService) CountyByAreaId(ctx context.Context, areaId string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".CountyByAreaId")
	defer span.End()

	area, err := ads.areaDao.CountyByAreaId(spanCtx, areaId, "")
	if err != nil {
		return nil, err
	}
	return &commonDo.Area{
		AreaId:       area.AreaId,
		AreaName:     area.AreaName,
		ParentAreaId: area.ParentId,
	}, nil
}

// ProvinceByAreaName 根据省地区名称查询省地区信息
// @param areaName 省地区名称
func (ads *areaDomainService) ProvinceByAreaName(ctx context.Context, areaName string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".ProvinceByAreaName")
	defer span.End()

	area, err := ads.areaDao.ProvinceByAreaName(spanCtx, areaName)
	if err != nil {
		return nil, err
	}
	return &commonDo.Area{
		AreaId:   area.AreaId,
		AreaName: area.AreaName,
	}, nil
}

// CityByAreaName 根据市地区名称查询市地区信息
// @param areaName 市地区名称
func (ads *areaDomainService) CityByAreaName(ctx context.Context, areaName string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".CityByAreaName")
	defer span.End()

	area, err := ads.areaDao.CityByAreaName(spanCtx, areaName)
	if err != nil {
		return nil, err
	}
	return &commonDo.Area{
		AreaId:       area.AreaId,
		AreaName:     area.AreaName,
		ParentAreaId: area.ParentId,
	}, nil
}

// CountyByAreaName 根据区县地区名称查询区县地区信息
// @param areaName 区县地区名称
func (ads *areaDomainService) CountyByAreaName(ctx context.Context, areaName string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".CountyByAreaName")
	defer span.End()

	area, err := ads.areaDao.CountyByAreaName(spanCtx, areaName)
	if err != nil {
		return nil, err
	}
	return &commonDo.Area{
		AreaId:       area.AreaId,
		AreaName:     area.AreaName,
		ParentAreaId: area.ParentId,
	}, nil
}

// CascadeArea 级联查询
// @param provinceAreaId 省地区编号
// @param cityAreaId 市地区编号
// @param countyAreaId 区县地区编号
func (ads *areaDomainService) CascadeArea(
	ctx context.Context,
	provinceAreaId string,
	cityAreaId string,
	countyAreaId string,
) (*commonDo.CascadeArea, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".CascadeArea")
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

	province, err := ads.areaDao.ProvinceByAreaId(spanCtx, provinceAreaId)
	if err != nil {
		return nil, err
	}
	if province.ID <= 0 {
		return nil, errors.New("省地区错误")
	}

	city, err := ads.areaDao.CityByAreaId(spanCtx, cityAreaId, province.AreaId)
	if err != nil {
		return nil, err
	}
	if city.ID <= 0 {
		return nil, errors.New("市地区错误")
	}

	county, err := ads.areaDao.CountyByAreaId(spanCtx, countyAreaId, city.AreaId)
	if err != nil {
		return nil, err
	}
	if county.ID <= 0 {
		return nil, errors.New("区县地区错误")
	}

	return &commonDo.CascadeArea{
		Province: commonDo.Area{AreaId: province.AreaId, AreaName: province.AreaName},
		City:     commonDo.Area{AreaId: city.AreaId, AreaName: city.AreaName, ParentAreaId: city.ParentId},
		County:   commonDo.Area{AreaId: county.AreaId, AreaName: county.AreaName, ParentAreaId: county.ParentId},
	}, nil
}

// Area 查询下级地区列表
// @param areaType 地区类型 省：province 市：city 区县：county
// @param parentAreaId 父级地区ID
func (ads *areaDomainService) Area(ctx context.Context, areaType string, parentAreaId string) ([]commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ads.baseTraceSpanName+".Area")
	defer span.End()

	if areaType == "" {
		return nil, errors.New("缺少地区类型")
	}
	if areaType != "province" && parentAreaId == "" {
		return nil, errors.New("缺少地区父ID")
	}

	var areaCollection []commonDo.Area

	switch areaType {
	case "province":
		provinceList, err := ads.areaDao.ProvinceList(spanCtx)
		if err != nil {
			logger.Error(spanCtx, "查询省列表失败", logger.ErrorField(err))
			return nil, errors.New("查询省列表失败")
		}
		for _, province := range provinceList {
			areaCollection = append(areaCollection, commonDo.Area{AreaId: province.AreaId, AreaName: province.AreaName})
		}

	case "city":
		cityList, err := ads.areaDao.CityList(spanCtx, parentAreaId)
		if err != nil {
			logger.Error(spanCtx, "查询市列表失败", logger.ErrorField(err))
			return nil, errors.New("查询市列表失败")
		}
		for _, city := range cityList {
			areaCollection = append(areaCollection, commonDo.Area{
				AreaId:       city.AreaId,
				AreaName:     city.AreaName,
				ParentAreaId: city.ParentId,
			})
		}

	case "county":
		countyList, err := ads.areaDao.CountyList(spanCtx, parentAreaId)
		if err != nil {
			logger.Error(spanCtx, "查询区县列表失败", logger.ErrorField(err))
			return nil, errors.New("查询区县列表失败")
		}
		for _, county := range countyList {
			areaCollection = append(areaCollection, commonDo.Area{
				AreaId:       county.AreaId,
				AreaName:     county.AreaName,
				ParentAreaId: county.ParentId,
			})
		}

	default:
		return nil, errors.New("地区类型错误")
	}

	return areaCollection, nil
}
