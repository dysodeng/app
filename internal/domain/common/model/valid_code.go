package model

// SenderType 消息发送者类型
type SenderType string

const (
	SmsSender   SenderType = "sms"
	EmailSender SenderType = "email"
)

// ValidCode 验证码结构
type ValidCode struct {
	Code   string `redis:"code"`
	Expire int64  `redis:"expire"`
	Time   int64  `redis:"time"`
}
