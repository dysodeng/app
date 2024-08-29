package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Array[T comparable] []T

func (a Array[T]) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Array[T]) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), a)
}

func (a *Array[T]) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}
