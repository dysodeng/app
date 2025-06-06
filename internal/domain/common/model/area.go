package model

import "github.com/dysodeng/app/internal/infrastructure/persistence/model/common"

// CascadeArea 级联地区
type CascadeArea struct {
	Province Province `json:"province"`
	City     City     `json:"city"`
	County   County   `json:"county"`
}

func CascadeAreaFromAreas(province Province, city City, county County) *CascadeArea {
	return &CascadeArea{
		Province: province,
		City:     city,
		County:   county,
	}
}

// Area 地区通用数据结构
type Area struct {
	AreaId       string `json:"area_id"`
	AreaName     string `json:"area_name"`
	ParentAreaId string `json:"parent_area_id"`
}

func AreaFromProvince(area Province) *Area {
	return &Area{
		AreaId:   area.AreaId,
		AreaName: area.AreaName,
	}
}

func AreaListFromProvince(areas []Province) []Area {
	result := make([]Area, len(areas))
	for i, area := range areas {
		result[i] = *AreaFromProvince(area)
	}
	return result
}

func AreaFromCity(area City) *Area {
	return &Area{
		AreaId:       area.AreaId,
		AreaName:     area.AreaName,
		ParentAreaId: area.ParentId,
	}
}

func AreaListFromCity(areas []City) []Area {
	result := make([]Area, len(areas))
	for i, area := range areas {
		result[i] = *AreaFromCity(area)
	}
	return result
}

func AreaFromCounty(area County) *Area {
	return &Area{
		AreaId:       area.AreaId,
		AreaName:     area.AreaName,
		ParentAreaId: area.ParentId,
	}
}

func AreaListFromCounty(areas []County) []Area {
	result := make([]Area, len(areas))
	for i, area := range areas {
		result[i] = *AreaFromCounty(area)
	}
	return result
}

type BigArea struct {
	ID       uint64     `json:"id"`
	AreaName string     `json:"area_name"`
	Province []Province `json:"province"`
}

type Province struct {
	ID            uint64 `json:"id"`
	AreaName      string `json:"area_name"`
	ShortAreaName string `json:"short_area_name"`
	AreaId        string `json:"area_id"`
	BigAreaID     uint64 `json:"big_area_id"`
}

func ProvinceFromModel(area common.Province) *Province {
	return &Province{
		ID:            area.ID,
		AreaName:      area.AreaName,
		ShortAreaName: area.ShortAreaName,
		AreaId:        area.AreaId,
		BigAreaID:     area.BigAreaID,
	}
}

func ProvinceListFromModel(areas []common.Province) []Province {
	result := make([]Province, len(areas))
	for i, area := range areas {
		result[i] = *ProvinceFromModel(area)
	}
	return result
}

type City struct {
	ID            uint64 `json:"id"`
	AreaName      string `json:"area_name"`
	AreaId        string `json:"area_id"`
	ShortAreaName string `json:"short_area_name"`
	ParentId      string `json:"parent_id"`
	ParentName    string `json:"parent_name"`
	IsOpen        uint8  `json:"is_open"`
	IsHot         uint8  `json:"is_hot"`
}

func CityFromModel(area common.City) *City {
	return &City{
		ID:            area.ID,
		AreaName:      area.AreaName,
		AreaId:        area.AreaId,
		ShortAreaName: area.ShortAreaName,
		ParentId:      area.ParentId,
		ParentName:    area.ParentName,
		IsOpen:        area.IsOpen,
		IsHot:         area.IsHot,
	}
}

func CityListFromModel(areas []common.City) []City {
	result := make([]City, len(areas))
	for i, area := range areas {
		result[i] = *CityFromModel(area)
	}
	return result
}

type County struct {
	ID         uint64 `json:"id"`
	AreaName   string `json:"area_name"`
	AreaId     string `json:"area_id"`
	ParentId   string `json:"parent_id"`
	ParentName string `json:"parent_name"`
}

func CountyFromModel(area common.County) *County {
	return &County{
		ID:         area.ID,
		AreaName:   area.AreaName,
		AreaId:     area.AreaId,
		ParentId:   area.ParentId,
		ParentName: area.ParentName,
	}
}

func CountyListFromModel(areas []common.County) []County {
	result := make([]County, len(areas))
	for i, county := range areas {
		result[i] = *CountyFromModel(county)
	}
	return result
}
