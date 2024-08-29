package area

import (
	"github.com/dysodeng/app/internal/model/common"
	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/service"
	"github.com/pkg/errors"
)

// Service 地区服务
type Service struct{}

func NewAreaService() *Service {
	return &Service{}
}

// ProvinceByAreaId 获取省地区信息
func (Service) ProvinceByAreaId(areaId string) (*common.Province, service.Error) {
	if areaId == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区编号")}
	}

	var province common.Province
	db.DB().Where("province_area_id=?", areaId).First(&province)
	if province.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("地区不存在")}
	}

	return &province, service.Error{Code: api.CodeOk}
}

// CityByAreaId 获取市地区信息
func (Service) CityByAreaId(areaId string) (*common.City, service.Error) {
	if areaId == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区编号")}
	}

	var city common.City
	db.DB().Where("city_area_id=?", areaId).First(&city)
	if city.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("地区不存在")}
	}

	return &city, service.Error{Code: api.CodeOk}
}

// CountyByAreaId 获取区县地区信息
func (Service) CountyByAreaId(areaId string) (*common.County, service.Error) {
	if areaId == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区编号")}
	}

	var county common.County
	db.DB().Where("county_area_id=?", areaId).First(&county)
	if county.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("地区不存在")}
	}

	return &county, service.Error{Code: api.CodeOk}
}

// ProvinceByAreaName 根据省名称获取地区信息
func (Service) ProvinceByAreaName(areaName string) (*common.Province, service.Error) {
	if areaName == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区名称")}
	}

	var province common.Province
	db.DB().Where("province_name=?", areaName).First(&province)
	if province.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("地区不存在")}
	}

	return &province, service.Error{Code: api.CodeOk}
}

// CityByAreaName 根据市名称获取地区信息
func (Service) CityByAreaName(areaName string) (*common.City, service.Error) {
	if areaName == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区名称")}
	}

	var city common.City
	db.DB().Where("city_name=?", areaName).First(&city)
	if city.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("地区不存在")}
	}

	return &city, service.Error{Code: api.CodeOk}
}

// CountyByAreaName 根据区县名称获取地区信息
func (Service) CountyByAreaName(areaName string) (*common.County, service.Error) {
	if areaName == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区名称")}
	}

	var county common.County
	db.DB().Where("county_name=?", areaName).First(&county)
	if county.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("地区不存在")}
	}

	return &county, service.Error{Code: api.CodeOk}
}

// CascadeArea 获取省市区级联地区信息
func (Service) CascadeArea(
	provinceAreaId string,
	cityAreaId string,
	countyAreaId string,
) (*common.CascadeArea, service.Error) {
	if provinceAreaId == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("省地区编号为空")}
	}
	if cityAreaId == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("市地区编号为空")}
	}
	if countyAreaId == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("区县地区编号为空")}
	}

	var province common.Province
	db.DB().Where("province_area_id=?", provinceAreaId).First(&province)
	if province.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("省地区错误")}
	}

	var city common.City
	db.DB().Where("city_area_id=?", cityAreaId).Where("city_parent_id=?", province.AreaId).First(&city)
	if city.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("市地区错误")}
	}

	var county common.County
	db.DB().Where("county_area_id=?", countyAreaId).Where("county_parent_id=?", city.AreaId).First(&county)
	if county.ID <= 0 {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("区县地区错误")}
	}

	return &common.CascadeArea{
		Province: province,
		City:     city,
		County:   county,
	}, service.Error{Code: api.CodeOk}
}

// Area 根据父级地区ID获取下级地区列表
func (Service) Area(areaType string, parentAreaId string) ([]common.Area, service.Error) {
	if areaType == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区类型")}
	}
	if areaType != "province" && parentAreaId == "" {
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少地区父ID")}
	}

	var areaCollection []common.Area
	switch areaType {
	case "province": // 省
		var provinceList []common.Province
		db.DB().Find(&provinceList)
		if len(provinceList) > 0 {
			for _, province := range provinceList {
				areaCollection = append(areaCollection, common.Area{
					AreaId:   province.AreaId,
					AreaName: province.AreaName,
				})
			}
		}
		break

	case "city": // 市
		var cityList []common.City
		db.DB().Where("city_parent_id=?", parentAreaId).Find(&cityList)
		if len(cityList) > 0 {
			for _, city := range cityList {
				areaCollection = append(areaCollection, common.Area{
					AreaId:       city.AreaId,
					AreaName:     city.AreaName,
					ParentAreaId: city.ParentId,
				})
			}
		}
		break

	case "county": // 区县
		var countyList []common.County
		db.DB().Where("county_parent_id=?", parentAreaId).Find(&countyList)
		if len(countyList) > 0 {
			for _, county := range countyList {
				areaCollection = append(areaCollection, common.Area{
					AreaId:       county.AreaId,
					AreaName:     county.AreaName,
					ParentAreaId: county.ParentId,
				})
			}
		}
		break

	default:
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("地区类型错误")}
	}

	return areaCollection, service.Error{Code: api.CodeOk}
}
