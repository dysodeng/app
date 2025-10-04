package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// CstHour 东八区
const CstHour int64 = 8 * 3600

const (
	TimeCompactLayout = "20060102150405"
	TimeDotLayout     = "2006.01.02 15:04:05"
	DateDotLayout     = "2006.01.02"
	DateCompactLayout = "20060102"
	ChineseDateLayout = "2006年01月02日"
	MonthLayout       = "2006-01"
)

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	zero, _ := time.Parse(time.DateTime, "0001-01-01 00:00:00")
	zeroTime := JSONTime{Time: zero}
	if t == zeroTime {
		return []byte(fmt.Sprintf("\"%s\"", "")), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(time.DateTime))
	return []byte(formatted), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		return
	}

	// 指定解析的格式
	now, err := time.Parse(time.TimeOnly, strings.Trim(string(data), "\""))
	*t = JSONTime{Time: now}
	return
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueOf time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// JSONDate format json date field by myself
type JSONDate struct {
	time.Time
}

func (t JSONDate) MarshalJSON() ([]byte, error) {
	zero, _ := time.Parse(time.DateOnly, "0001-01-01")
	zeroTime := JSONDate{Time: zero}
	if t == zeroTime {
		return []byte(fmt.Sprintf("\"%s\"", "")), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(time.DateOnly))
	return []byte(formatted), nil
}

func (t *JSONDate) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		return
	}

	// 指定解析的格式
	now, err := time.Parse(time.DateOnly, strings.Trim(string(data), "\""))
	*t = JSONDate{Time: now}
	return
}

// Value insert timestamp into mysql need this function.
func (t JSONDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueOf time.Time
func (t *JSONDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// TimeOnly 只包含时间的类型
type TimeOnly struct {
	time.Time
}

// NewTimeOnly 创建 TimeOnly 实例
func NewTimeOnly(hour, minute, second int) TimeOnly {
	return TimeOnly{
		Time: time.Date(0, 1, 1, hour, minute, second, 0, time.UTC),
	}
}

// ParseTimeOnly 从字符串解析时间
func ParseTimeOnly(timeStr string) (TimeOnly, error) {
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return TimeOnly{}, err
	}
	return TimeOnly{Time: t}, nil
}

// String 返回时间字符串格式
func (t TimeOnly) String() string {
	return t.Format("15:04:05")
}

// MarshalJSON JSON 序列化
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON JSON 反序列化
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}
	str = str[1 : len(str)-1] // 去掉引号
	parsed, err := ParseTimeOnly(str)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

// Value 实现 driver.Valuer 接口
func (t TimeOnly) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan 实现 sql.Scanner 接口
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		parsed, err := ParseTimeOnly(v)
		if err != nil {
			return err
		}
		*t = parsed
	case []byte:
		parsed, err := ParseTimeOnly(string(v))
		if err != nil {
			return err
		}
		*t = parsed
	case time.Time:
		*t = TimeOnly{Time: v}
	default:
		return fmt.Errorf("cannot scan %T into TimeOnly", value)
	}
	return nil
}
