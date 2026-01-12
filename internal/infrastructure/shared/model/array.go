package model

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
)

type Array[T comparable] []T

func (a Array[T]) Value() (driver.Value, error) {
	return sonic.Marshal(a)
}

func (a *Array[T]) Scan(v interface{}) error {
	return sonic.Unmarshal(v.([]byte), a)
}

func (a *Array[T]) String() string {
	b, _ := sonic.Marshal(a)
	return string(b)
}
