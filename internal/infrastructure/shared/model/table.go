package model

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PrimaryKeyID 自增主键ID
type PrimaryKeyID struct {
	ID uint64 `gorm:"primary_key;autoIncrement" json:"id"`
}

// DistributedPrimaryKeyID 分布式主键ID
type DistributedPrimaryKeyID struct {
	ID uuid.UUID `gorm:"type:uuid;not null;default:uuid_generate_v7();primary_key" json:"id"`
}

// Time 添加时间,修改时间
type Time struct {
	CreatedAt JSONTime `gorm:"type:timestamp(0) without time zone;index;not null" json:"created_at,omitempty"`
	UpdatedAt JSONTime `gorm:"type:timestamp(0) without time zone;not null" json:"updated_at,omitempty"`
}

// SoftDelete 软删除
type SoftDelete struct {
	IsDelete   uint8    `gorm:"not null;default:0;comment:删除标识 0-未删除 1-已删除" json:"is_delete"`
	DeleteTime JSONTime `gorm:"type:timestamp(0) without time zone;index;not null" json:"delete_time"`
}

// TableComment 修改表注释
func TableComment(tx *gorm.DB, driver, tableName, comment string) {
	switch driver {
	case "mysql":
		tx.Exec(fmt.Sprintf(`ALTER TABLE %s COMMENT="%s"`, tableName, comment))
	case "postgres":
		tx.Exec(fmt.Sprintf(`COMMENT ON TABLE "%s" IS '%s'`, tableName, comment))
	}
}
