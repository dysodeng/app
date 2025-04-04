package websocket

import (
	"log"

	"github.com/dysodeng/app/internal/pkg/ws"
	wsMessage "github.com/dysodeng/app/internal/pkg/ws/message"
)

type MessageHandler struct{}

// Handler 处理消息
func (MessageHandler) Handler(clientId, userId, userType string, messageType int, message []byte) error {
	log.Println("message type: ", messageType)
	log.Printf("receive message: %s", message)
	_ = ws.SendMessage(clientId, wsMessage.TypeMessage, "hello world")
	return nil
}
