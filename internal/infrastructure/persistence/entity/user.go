package entity

import (
	"time"

	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Username  string         `gorm:"size:50;not null;uniqueIndex"`
	Email     string         `gorm:"size:100;not null;uniqueIndex"`
	Password  string         `gorm:"size:100;not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
