package user

import "github.com/dysodeng/app/internal/pkg/model"

// User 用户表
type User struct {
	model.PrimaryKeyID
	Telephone string             `gorm:"column:telephone;type:char(11);index:user_telephone_idx,unique;not null;comment:手机号码" json:"telephone"`
	Password  string             `gorm:"column:password;type:varchar(60);not null;default:'';comment:登录密码" json:"-"`
	RealName  string             `gorm:"column:real_name;type:varchar(20);not null;default:'';comment:真实姓名" json:"real_name"`
	Avatar    string             `gorm:"column:avatar;type:varchar(150);not null;default:'';comment:头像" json:"avatar"`
	Nickname  string             `gorm:"column:nickname;type:varchar(20);not null;default:'';comment:昵称" json:"nickname"`
	Gender    uint8              `gorm:"column:gender;type:tinyint;not null;default:0;comment:性别 0-保密 1-男 2-女" json:"gender"`
	Birthday  model.JSONDate     `gorm:"column:birthday;type:date;comment:生日" json:"birthday"`
	Status    model.BinaryStatus `gorm:"column:status;type:tinyint;not null;default:1;comment:状态 0-禁用 1-启用" json:"status"`
	model.Time
}

func (User) TableName() string {
	return "user"
}
