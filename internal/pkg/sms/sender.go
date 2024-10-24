package sms

type Sender interface {
	// SendSms 发送短信
	SendSms() (bool, error)
}
