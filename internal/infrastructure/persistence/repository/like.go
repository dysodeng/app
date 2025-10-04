package repository

import (
	"fmt"

	"gorm.io/gorm"
)

func WhereLike(tx *gorm.DB, field, value string) *gorm.DB {
	return tx.Where(field+" LIKE ?", fmt.Sprintf("%%%%%s%%%%", value))
}

func WhereRightLike(tx *gorm.DB, field, value string) *gorm.DB {
	return tx.Where(field+" LIKE ?", fmt.Sprintf("%s%%%%", value))
}
