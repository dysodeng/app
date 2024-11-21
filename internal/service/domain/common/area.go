package common

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/telemetry/trace"

	commonDao "github.com/dysodeng/app/internal/dal/dao/common"
	"github.com/dysodeng/app/internal/pkg/logger"
	commonDo "github.com/dysodeng/app/internal/service/do/common"
	"github.com/pkg/errors"
)

// AreaDomainService Area领域服务
type AreaDomainService struct {
	ctx               context.Context
	baseTraceSpanName string
}

func NewAreaDomainService(ctx context.Context) *AreaDomainService {
	return &AreaDomainService{
		ctx:               ctx,
		baseTraceSpanName: "service.domain.common.AreaDomainService",
	}
}

// ProvinceByAreaId 根据省地区编号查询省地区信息
// @param areaId 省地区编号
func (ads *AreaDomainService) ProvinceByAreaId(areaId string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".ProvinceByAreaId")
	defer span.End()

	areaDao := commonDao.NewAreaDao(spanCtx)
	area, err := areaDao.ProvinceByAreaId(areaId)
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
func (ads *AreaDomainService) CityByAreaId(areaId string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".CityByAreaId")
	defer span.End()

	areaDao := commonDao.NewAreaDao(spanCtx)
	area, err := areaDao.CityByAreaId(areaId, "")
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
func (ads *AreaDomainService) CountyByAreaId(areaId string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".CountyByAreaId")
	defer span.End()

	areaDao := commonDao.NewAreaDao(spanCtx)
	area, err := areaDao.CountyByAreaId(areaId, "")
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
func (ads *AreaDomainService) ProvinceByAreaName(areaName string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".ProvinceByAreaName")
	defer span.End()

	areaDao := commonDao.NewAreaDao(spanCtx)
	area, err := areaDao.ProvinceByAreaName(areaName)
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
func (ads *AreaDomainService) CityByAreaName(areaName string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".CityByAreaName")
	defer span.End()

	areaDao := commonDao.NewAreaDao(spanCtx)
	area, err := areaDao.CityByAreaName(areaName)
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
func (ads *AreaDomainService) CountyByAreaName(areaName string) (*commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".CountyByAreaName")
	defer span.End()

	areaDao := commonDao.NewAreaDao(spanCtx)
	area, err := areaDao.CountyByAreaName(areaName)
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
func (ads *AreaDomainService) CascadeArea(
	provinceAreaId string,
	cityAreaId string,
	countyAreaId string,
) (*commonDo.CascadeArea, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".CascadeArea")
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

	areaDao := commonDao.NewAreaDao(spanCtx)
	province, err := areaDao.ProvinceByAreaId(provinceAreaId)
	if err != nil {
		return nil, err
	}
	if province.ID <= 0 {
		return nil, errors.New("省地区错误")
	}

	city, err := areaDao.CityByAreaId(cityAreaId, province.AreaId)
	if err != nil {
		return nil, err
	}
	if city.ID <= 0 {
		return nil, errors.New("市地区错误")
	}

	county, err := areaDao.CountyByAreaId(countyAreaId, city.AreaId)
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
func (ads *AreaDomainService) Area(areaType string, parentAreaId string) ([]commonDo.Area, error) {
	spanCtx, span := trace.Tracer().Start(ads.ctx, ads.baseTraceSpanName+".Area")
	defer span.End()

	if areaType == "" {
		return nil, errors.New("缺少地区类型")
	}
	if areaType != "province" && parentAreaId == "" {
		return nil, errors.New("缺少地区父ID")
	}

	var areaCollection []commonDo.Area

	areaDao := commonDao.NewAreaDao(spanCtx)

	switch areaType {
	case "province":
		provinceList, err := areaDao.ProvinceList()
		if err != nil {
			logger.Error(spanCtx, "查询省列表失败", logger.ErrorField(err))
			return nil, errors.New("查询省列表失败")
		}
		for _, province := range provinceList {
			areaCollection = append(areaCollection, commonDo.Area{AreaId: province.AreaId, AreaName: province.AreaName})
		}
		break

	case "city":
		cityList, err := areaDao.CityList(parentAreaId)
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
		break

	case "county":
		countyList, err := areaDao.CountyList(parentAreaId)
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
		break

	default:
		return nil, errors.New("地区类型错误")
	}

	return areaCollection, nil
}
