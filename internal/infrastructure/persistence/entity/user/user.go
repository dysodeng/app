package user

import "github.com/dysodeng/app/internal/infrastructure/shared/model"

type User struct {
	model.DistributedPrimaryKeyID
	Telephone           string `gorm:"type:varchar(15);index:user_telephone_idx,unique;not null;default:'';comment:手机号" json:"telephone"`
	WxUnionID           string `gorm:"type:varchar(36);index:user_wx_union_idx;not null;default:'';comment:微信开放平台用户UnionID" json:"wx_union_id"`
	WxMiniProgramOpenID string `gorm:"column:wx_mini_program_openid;type:varchar(36);index:user_wx_mp_idx,unique;not null;default:'';comment:微信小程序用户OpenID" json:"wx_mini_program_openid"`
	WxOfficialOpenID    string `gorm:"column:wx_official_openid;type:varchar(36);index:user_wx_official_idx;not null;default:'';comment:微信公众号用户OpenID" json:"wx_official_openid"`
	Nickname            string `gorm:"type:varchar(50);not null;default:'';comment:用户昵称" json:"nickname"`
	Avatar              string `gorm:"type:varchar(150);not null;default:'';comment:用户头像" json:"avatar"`
	Status              uint8  `gorm:"not null;default:0;comment:状态 0-禁用 1-启用" json:"status"`
	model.Time
}

func (User) TableName() string {
	return "users"
}
