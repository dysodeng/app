package area

import (
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/service"
)

// ServiceInterface 地区服务
type ServiceInterface interface {
	// ProvinceByAreaId 获取省地区信息
	ProvinceByAreaId(areaId string) (*common.Province, service.Error)
	// CityByAreaId 获取市地区信息
	CityByAreaId(areaId string) (*common.City, service.Error)
	// CountyByAreaId 获取区县地区信息
	CountyByAreaId(areaId string) (*common.County, service.Error)
	// ProvinceByAreaName 根据省名称获取地区信息
	ProvinceByAreaName(areaName string) (*common.Province, service.Error)
	// CityByAreaName 根据市名称获取地区信息
	CityByAreaName(areaName string) (*common.City, service.Error)
	// CountyByAreaName 根据区县名称获取地区信息
	CountyByAreaName(areaName string) (*common.County, service.Error)
	// CascadeArea 获取省市区级联地区信息
	CascadeArea(provinceAreaId string, cityAreaId string, countyAreaId string) (*common.CascadeArea, service.Error)
	// Area 根据父级地区ID获取下级地区列表
	Area(areaType string, parentAreaId string) ([]common.Area, service.Error)
}
