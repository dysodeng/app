package common

// Area 地区领域对象
type Area struct {
	AreaId       string `json:"area_id"`
	AreaName     string `json:"area_name"`
	ParentAreaId string `json:"parent_area_id"`
}

// CascadeArea 级联地区领域对象
type CascadeArea struct {
	Province Area `json:"province"`
	City     Area `json:"city"`
	County   Area `json:"county"`
}
