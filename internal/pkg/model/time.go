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
	if t.Time.UnixNano() == zeroTime.UnixNano() {
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
	if t.Time.UnixNano() == zeroTime.UnixNano() {
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
