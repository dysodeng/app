package common

import (
	"github.com/dysodeng/app/internal/pkg/model"
)

// SenderType 消息发送者类型
type SenderType string

const (
	SmsSender   SenderType = "sms"
	EmailSender SenderType = "email"
)

// CodeTemplateType 验证码模板类型
type CodeTemplateType string

func (t CodeTemplateType) String() string {
	return string(t)
}

const (
	CodeLogin          CodeTemplateType = "login"           // 登录
	CodeForgetPassword CodeTemplateType = "forget_password" // 找回密码
	CodeBindAccount    CodeTemplateType = "bind_account"    // 绑定账号
	CodeWithdraw       CodeTemplateType = "withdraw"        // 提现
	Code               CodeTemplateType = "code"            // 常规验证码
)

// SmsConfig 短信配置
type SmsConfig struct {
	model.PrimaryKeyID
	SmsType         string `gorm:"type:varchar(150); not null;default:'';comment:短信服务商类型 ali_cloud-阿里云" json:"sms_type"`
	AppKey          string `gorm:"type:varchar(150); not null;default:'';comment:短信AppKey" json:"app_key"`
	SecretKey       string `gorm:"type:varchar(150); not null;default:'';comment:短信SecretKey" json:"secret_key"`
	FreeSignName    string `gorm:"type:varchar(150); not null;default:'';comment:短信签名" json:"free_sign_name"`
	ValidCodeExpire uint   `gorm:"type:int(10); not null;default:0;comment:短信验证码过期时间，单位分钟" json:"valid_code_expire"`
	model.Time
}

func (SmsConfig) TableName() string {
	return "sms_config"
}

// SmsTemplate 短信模版
type SmsTemplate struct {
	model.PrimaryKeyID
	TemplateName string `gorm:"type:varchar(150); not null; default:'';comment:模版名称" json:"template_name"`
	Template     string `gorm:"type:varchar(150); not null; default:'';comment:短信模版类型" json:"template"`
	TemplateId   string `gorm:"type:varchar(150); not null; default:'';comment:短信模版ID" json:"template_id"`
	model.Time
}

func (SmsTemplate) TableName() string {
	return "sms_template"
}
