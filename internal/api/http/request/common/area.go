package common

type AreaBody struct {
	AreaType     string `json:"area_type"`
	ParentAreaId string `json:"parent_area_id"`
}

type CascadeAreaBody struct {
	ProvinceAreaId string `json:"province_area_id"`
	CityAreaId     string `json:"city_area_id"`
	CountyAreaId   string `json:"county_area_id"`
}
