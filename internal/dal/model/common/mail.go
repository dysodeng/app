package common

import (
	"github.com/dysodeng/app/internal/pkg/model"
)

// MailConfig 邮件配置
type MailConfig struct {
	model.PrimaryKeyID
	User      string `gorm:"type:varchar(150);not null;default:'';comment:发件邮箱地址" json:"user"`
	FromName  string `gorm:"type:varchar(150);not null;default:'';comment:发送者名称" json:"from_name"`
	Transport string `gorm:"type:varchar(150);not null;default:'smtp';comment:邮件传输协议" json:"transport"`
	Username  string `gorm:"type:varchar(150);not null;default:'';comment:用户名" json:"username"`
	Password  string `gorm:"type:varchar(150);not null;default:'';comment:密码" json:"password"`
	Host      string `gorm:"type:varchar(150);not null;default:'';comment:邮件服务器地址" json:"host"`
	Port      int    `gorm:"not null;default:0;comment:邮件服务器端口" json:"port"`
	model.Time
}

func (MailConfig) TableName() string {
	return "mail_config"
}
