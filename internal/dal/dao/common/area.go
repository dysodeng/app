package common

import (
	"context"

	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// AreaDao area数据操作层
type AreaDao struct {
	ctx context.Context
}

func NewAreaDao(ctx context.Context) *AreaDao {
	return &AreaDao{ctx: ctx}
}

func (dao *AreaDao) ProvinceByAreaId(areaId string) (*common.Province, error) {
	var province common.Province
	err := db.DB().WithContext(dao.ctx).Where("province_area_id=?", areaId).First(&province).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &province, nil
}

func (dao *AreaDao) CityByAreaId(areaId, parentAreaId string) (*common.City, error) {
	var city common.City
	query := db.DB().WithContext(dao.ctx).Where("city_area_id=?", areaId)
	if parentAreaId != "" {
		query = query.Where("city_parent_id=?", parentAreaId)
	}
	err := query.First(&city).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &city, nil
}

func (dao *AreaDao) CountyByAreaId(areaId, parentAreaId string) (*common.County, error) {
	var county common.County
	query := db.DB().WithContext(dao.ctx).Where("county_area_id=?", areaId)
	if parentAreaId != "" {
		query = query.Where("county_parent_id=?", parentAreaId)
	}
	err := query.First(&county).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &county, nil
}

func (dao *AreaDao) ProvinceByAreaName(areaName string) (*common.Province, error) {
	var province common.Province
	err := db.DB().WithContext(dao.ctx).Where("province_name=?", areaName).First(&province).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &province, nil
}

func (dao *AreaDao) CityByAreaName(areaName string) (*common.City, error) {
	var city common.City
	err := db.DB().WithContext(dao.ctx).Where("city_name=?", areaName).First(&city).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &city, nil
}

func (dao *AreaDao) CountyByAreaName(areaName string) (*common.County, error) {
	var county common.County
	err := db.DB().WithContext(dao.ctx).Where("county_name=?", areaName).First(&county).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &county, nil
}

func (dao *AreaDao) ProvinceList() ([]common.Province, error) {
	var provinceList []common.Province
	err := db.DB().WithContext(dao.ctx).Find(&provinceList).Error
	if err != nil {
		return nil, err
	}
	return provinceList, nil
}

func (dao *AreaDao) CityList(provinceAreaId string) ([]common.City, error) {
	var cityList []common.City
	err := db.DB().WithContext(dao.ctx).Where("city_parent_id=?", provinceAreaId).Find(&cityList).Error
	if err != nil {
		return nil, err
	}
	return cityList, nil
}

func (dao *AreaDao) CountyList(cityAreaId string) ([]common.County, error) {
	var countyList []common.County
	err := db.DB().WithContext(dao.ctx).Where("county_parent_id=?", cityAreaId).Find(&countyList).Error
	if err != nil {
		return nil, err
	}
	return countyList, nil
}
