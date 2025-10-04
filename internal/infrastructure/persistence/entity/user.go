package entity

import (
	"github.com/dysodeng/app/internal/infrastructure/shared/model"
)

// User 用户实体
type User struct {
	model.DistributedPrimaryKeyID
	Username string `gorm:"size:50;not null;uniqueIndex"`
	Email    string `gorm:"size:100;not null;uniqueIndex"`
	Password string `gorm:"size:100;not null"`
	model.Time
}

func (User) TableName() string {
	return "users"
}
