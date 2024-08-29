package model

import (
	"github.com/dysodeng/app/internal/pkg/form"
	"gorm.io/gorm"
)

// Pagination 分页排序
func Pagination(pagination form.Pagination, order ...string) func(tx *gorm.DB) *gorm.DB {
	pagination.CheckOrDefault()
	return func(tx *gorm.DB) *gorm.DB {
		if len(order) > 0 {
			for _, s := range order {
				tx.Order(s)
			}
		}
		return tx.Offset(pagination.Offset()).Limit(pagination.PageSize)
	}
}
