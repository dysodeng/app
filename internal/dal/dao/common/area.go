package common

import (
	"context"

	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// AreaDao 地区数据访问层
type AreaDao interface {
	ProvinceByAreaId(ctx context.Context, areaId string) (*common.Province, error)
	ProvinceByAreaName(ctx context.Context, areaName string) (*common.Province, error)
	CityByAreaId(ctx context.Context, areaId string, parentAreaId string) (*common.City, error)
	CityByAreaName(ctx context.Context, areaName string) (*common.City, error)
	CountyByAreaId(ctx context.Context, areaId string, parentAreaId string) (*common.County, error)
	CountyByAreaName(ctx context.Context, areaName string) (*common.County, error)
	ProvinceList(ctx context.Context) ([]common.Province, error)
	CityList(ctx context.Context, provinceAreaId string) ([]common.City, error)
	CountyList(ctx context.Context, cityAreaId string) ([]common.County, error)
}

// areaDao area数据操作层
type areaDao struct {
	baseTraceSpanName string
}

func NewAreaDao() AreaDao {
	return &areaDao{
		baseTraceSpanName: "dal.dao.common.AreaDao",
	}
}

func (dao *areaDao) ProvinceByAreaId(ctx context.Context, areaId string) (*common.Province, error) {
	var province common.Province
	err := db.DB().WithContext(ctx).Where("province_area_id=?", areaId).First(&province).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &province, nil
}

func (dao *areaDao) CityByAreaId(ctx context.Context, areaId, parentAreaId string) (*common.City, error) {
	var city common.City
	query := db.DB().WithContext(ctx).Where("city_area_id=?", areaId)
	if parentAreaId != "" {
		query = query.Where("city_parent_id=?", parentAreaId)
	}
	err := query.First(&city).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &city, nil
}

func (dao *areaDao) CountyByAreaId(ctx context.Context, areaId, parentAreaId string) (*common.County, error) {
	var county common.County
	query := db.DB().WithContext(ctx).Where("county_area_id=?", areaId)
	if parentAreaId != "" {
		query = query.Where("county_parent_id=?", parentAreaId)
	}
	err := query.First(&county).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &county, nil
}

func (dao *areaDao) ProvinceByAreaName(ctx context.Context, areaName string) (*common.Province, error) {
	var province common.Province
	err := db.DB().WithContext(ctx).Where("province_name=?", areaName).First(&province).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &province, nil
}

func (dao *areaDao) CityByAreaName(ctx context.Context, areaName string) (*common.City, error) {
	var city common.City
	err := db.DB().WithContext(ctx).Where("city_name=?", areaName).First(&city).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &city, nil
}

func (dao *areaDao) CountyByAreaName(ctx context.Context, areaName string) (*common.County, error) {
	var county common.County
	err := db.DB().WithContext(ctx).Where("county_name=?", areaName).First(&county).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &county, nil
}

func (dao *areaDao) ProvinceList(ctx context.Context) ([]common.Province, error) {
	var provinceList []common.Province
	err := db.DB().WithContext(ctx).Find(&provinceList).Error
	if err != nil {
		return nil, err
	}
	return provinceList, nil
}

func (dao *areaDao) CityList(ctx context.Context, provinceAreaId string) ([]common.City, error) {
	var cityList []common.City
	err := db.DB().WithContext(ctx).Where("city_parent_id=?", provinceAreaId).Find(&cityList).Error
	if err != nil {
		return nil, err
	}
	return cityList, nil
}

func (dao *areaDao) CountyList(ctx context.Context, cityAreaId string) ([]common.County, error) {
	var countyList []common.County
	err := db.DB().WithContext(ctx).Where("county_parent_id=?", cityAreaId).Find(&countyList).Error
	if err != nil {
		return nil, err
	}
	return countyList, nil
}
