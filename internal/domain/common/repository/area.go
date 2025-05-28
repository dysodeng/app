package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/common/model"
)

// AreaRepository area仓储服务
type AreaRepository interface {
	ProvinceByAreaId(ctx context.Context, areaId string) (*model.Province, error)
	ProvinceByAreaName(ctx context.Context, areaName string) (*model.Province, error)
	CityByAreaId(ctx context.Context, areaId string, parentAreaId string) (*model.City, error)
	CityByAreaName(ctx context.Context, areaName string) (*model.City, error)
	CountyByAreaId(ctx context.Context, areaId string, parentAreaId string) (*model.County, error)
	CountyByAreaName(ctx context.Context, areaName string) (*model.County, error)
	ProvinceList(ctx context.Context) ([]model.Province, error)
	CityList(ctx context.Context, provinceAreaId string) ([]model.City, error)
	CountyList(ctx context.Context, cityAreaId string) ([]model.County, error)
}
