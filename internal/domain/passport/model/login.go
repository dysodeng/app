package model

import (
	"github.com/dysodeng/app/internal/domain/passport/valueobject"
	"github.com/google/uuid"
)

// AdminLoginInfo 管理员登录信息
type AdminLoginInfo struct {
	AdminID     uint64
	Username    string
	IsSuper     bool
	Permissions []string
}

// UserLoginInfo 用户登录信息
type UserLoginInfo struct {
	Registered   bool // 用户是否已注册
	PlatformType valueobject.PlatformType
	UserId       uuid.UUID
	Telephone    string
	Avatar       string
}
