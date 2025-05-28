package common

import (
	"context"
	"github.com/dysodeng/app/internal/domain/common/model"
	"github.com/dysodeng/app/internal/domain/common/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/model/common"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type areaRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewAreaRepository(txManager transactions.TransactionManager) repository.AreaRepository {
	return &areaRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.common.AreaRepository",
		txManager:         txManager,
	}
}

func (repo *areaRepository) ProvinceByAreaId(ctx context.Context, areaId string) (*model.Province, error) {
	var province common.Province
	err := db.DB().WithContext(ctx).Where("province_area_id=?", areaId).First(&province).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.ProvinceFromModel(province), nil
}

func (repo *areaRepository) CityByAreaId(ctx context.Context, areaId, parentAreaId string) (*model.City, error) {
	var city common.City
	query := db.DB().WithContext(ctx).Where("city_area_id=?", areaId)
	if parentAreaId != "" {
		query = query.Where("city_parent_id=?", parentAreaId)
	}
	err := query.First(&city).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.CityFromModel(city), nil
}

func (repo *areaRepository) CountyByAreaId(ctx context.Context, areaId, parentAreaId string) (*model.County, error) {
	var county common.County
	query := db.DB().WithContext(ctx).Where("county_area_id=?", areaId)
	if parentAreaId != "" {
		query = query.Where("county_parent_id=?", parentAreaId)
	}
	err := query.First(&county).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.CountyFromModel(county), nil
}

func (repo *areaRepository) ProvinceByAreaName(ctx context.Context, areaName string) (*model.Province, error) {
	var province common.Province
	err := db.DB().WithContext(ctx).Where("province_name=?", areaName).First(&province).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.ProvinceFromModel(province), nil
}

func (repo *areaRepository) CityByAreaName(ctx context.Context, areaName string) (*model.City, error) {
	var city common.City
	err := db.DB().WithContext(ctx).Where("city_name=?", areaName).First(&city).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.CityFromModel(city), nil
}

func (repo *areaRepository) CountyByAreaName(ctx context.Context, areaName string) (*model.County, error) {
	var county common.County
	err := db.DB().WithContext(ctx).Where("county_name=?", areaName).First(&county).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.CountyFromModel(county), nil
}

func (repo *areaRepository) ProvinceList(ctx context.Context) ([]model.Province, error) {
	var provinceList []common.Province
	err := db.DB().WithContext(ctx).Find(&provinceList).Error
	if err != nil {
		return nil, err
	}
	return model.ProvinceListFromModel(provinceList), nil
}

func (repo *areaRepository) CityList(ctx context.Context, provinceAreaId string) ([]model.City, error) {
	var cityList []common.City
	err := db.DB().WithContext(ctx).Where("city_parent_id=?", provinceAreaId).Find(&cityList).Error
	if err != nil {
		return nil, err
	}
	return model.CityListFromModel(cityList), nil
}

func (repo *areaRepository) CountyList(ctx context.Context, cityAreaId string) ([]model.County, error) {
	var countyList []common.County
	err := db.DB().WithContext(ctx).Where("county_parent_id=?", cityAreaId).Find(&countyList).Error
	if err != nil {
		return nil, err
	}
	return model.CountyListFromModel(countyList), nil
}
