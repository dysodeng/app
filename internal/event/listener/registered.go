package listener

import (
	"log"
)

// Registered 用户注册成功事件
type Registered struct{}

func (Registered) Handle(data map[string]interface{}) {
	log.Println(data)
}
