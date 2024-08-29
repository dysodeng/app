package message

// Code 验证码结构
type Code struct {
	Code   string `redis:"code"`
	Expire int64  `redis:"expire"`
	Time   int64  `redis:"time"`
}
