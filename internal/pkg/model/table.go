package model

// PrimaryKeyID 自增主键ID
type PrimaryKeyID struct {
	ID uint64 `gorm:"primary_key;autoIncrement" json:"id"`
}

// Time 添加时间,修改时间
type Time struct {
	CreatedAt JSONTime `gorm:"type:datetime(0);index;not null" json:"created_at,omitempty"`
	UpdatedAt JSONTime `gorm:"type:datetime(0);not null" json:"updated_at,omitempty"`
}

// SoftDelete 软删除
type SoftDelete struct {
	IsDelete   uint8    `gorm:"not null;default:0;comment:删除标识 0-未删除 1-已删除" json:"is_delete"`
	DeleteTime JSONTime `gorm:"type:datetime(0);index;not null" json:"delete_time"`
}
