package model

import (
	"github.com/dysodeng/app/internal/domain/passport/valueobject"
	"github.com/google/uuid"
)

// UserLoginInfo 用户登录信息
type UserLoginInfo struct {
	Registered   bool // 用户是否已注册
	PlatformType valueobject.PlatformType
	UserId       uuid.UUID
	Telephone    string
	Avatar       string
}
