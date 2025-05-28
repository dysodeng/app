package link

import (
	"database/sql/driver"
	"encoding/json"
)

// Type 链接类型
type Type string

const (
	Empty       Type = "empty"        // 不设置跳转
	Navigation  Type = "navigation"   // 外部导航
	MiniProgram Type = "mini_program" // 小程序
)

var TypeList = []Type{
	Empty,
	Navigation,
	MiniProgram,
}

// Link 跳转链接
type Link struct {
	Type   Type   `json:"type"`
	Params Params `json:"params"`
}

// Params 链接参数
type Params struct {
	ID         *uint64 `json:"id,omitempty"`
	AppID      *string `json:"app_id,omitempty"`
	Path       *string `json:"path,omitempty"`
	EnvVersion *string `json:"env_version"`
	Point      *Point  `json:"point,omitempty"`
}

type Point struct {
	Name      *string  `json:"name,omitempty"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

func (l Link) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *Link) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), l)
}
