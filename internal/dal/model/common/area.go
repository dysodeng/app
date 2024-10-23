package common

import (
	"github.com/dysodeng/app/internal/pkg/model"
)

// CascadeArea 级联地区
type CascadeArea struct {
	Province Province `json:"province"`
	City     City     `json:"city"`
	County   County   `json:"county"`
}

// Area 地区通用数据结构
type Area struct {
	AreaId       string `json:"area_id"`
	AreaName     string `json:"area_name"`
	ParentAreaId string `json:"parent_area_id"`
}

// BigArea 大区
type BigArea struct {
	model.PrimaryKeyID
	AreaName string     `gorm:"column:big_area_name;type:varchar(100);not null;default:''" json:"area_name"`
	Province []Province `gorm:"foreignKey:BigAreaID;references:ID" json:"province"`
}

func (BigArea) TableName() string {
	return "address_area"
}

// Province 省
type Province struct {
	model.PrimaryKeyID
	AreaName      string `gorm:"column:province_name;type:varchar(100);not null;default:''" json:"area_name"`
	ShortAreaName string `gorm:"column:province_short_name;type:varchar(100);not null;default:''" json:"short_area_name"`
	AreaId        string `gorm:"column:province_area_id;type:varchar(100);not null;default:''" json:"area_id"`
	BigAreaID     uint64 `gorm:"not null;default:0" json:"big_area_id"`
	model.Time
}

func (Province) TableName() string {
	return "address_province"
}

// City 市
type City struct {
	model.PrimaryKeyID
	AreaName      string `gorm:"column:city_name;type:varchar(100);not null;default:''" json:"area_name"`
	AreaId        string `gorm:"column:city_area_id;type:varchar(100);not null;default:''" json:"area_id"`
	ShortAreaName string `gorm:"column:city_short_name;type:varchar(100);not null;default:''" json:"short_area_name"`
	ParentId      string `gorm:"column:city_parent_id;type:varchar(100);not null;default:''" json:"parent_id"`
	ParentName    string `gorm:"column:city_parent_name;type:varchar(100);not null;default:''" json:"parent_name"`
	IsOpen        uint8  `gorm:"not null;default:0;comment:是否已开通该城市 0-未开通 1-已开通" json:"is_open"`
	IsHot         uint8  `gorm:"not null;default:0;comment:是否热门城市 0-否 1-是" json:"is_hot"`
	model.Time
}

func (City) TableName() string {
	return "address_city"
}

// County 区县
type County struct {
	model.PrimaryKeyID
	AreaName   string `gorm:"column:county_name;type:varchar(100);not null;default:''" json:"area_name"`
	AreaId     string `gorm:"column:county_area_id;type:varchar(100);not null;default:''" json:"area_id"`
	ParentId   string `gorm:"column:county_parent_id;type:varchar(100);not null;default:''" json:"parent_id"`
	ParentName string `gorm:"column:county_parent_name;type:varchar(100);not null;default:''" json:"parent_name"`
	model.Time
}

func (County) TableName() string {
	return "address_county"
}
