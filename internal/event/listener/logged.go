package listener

import "log"

// Logged 用户登录成功事件
type Logged struct{}

func (handle Logged) Handle(data map[string]interface{}) {
	log.Println(data)
}
