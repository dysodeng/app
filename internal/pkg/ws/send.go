package ws

import (
	"fmt"

	"github.com/dysodeng/app/internal/pkg/ws/message"
	"github.com/pkg/errors"
)

func SendMessage(clientId string, msgType message.Type, msg string) error {
	// 从 hub 中获取 client_id 关联的 client
	HubBus.RLock()
	client, ok := HubBus.clients[clientId]
	HubBus.RUnlock()
	if !ok {
		return errors.New(fmt.Sprintf("client not found: %s", clientId))
	}

	// 消息装箱
	messageItem := message.WsMessage{
		ClientID: clientId,
		Message:  message.Message{Type: msgType, Data: msg},
	}

	// 发送消息
	client.send <- messageItem

	// 增加等待发送的消息数量
	HubBus.pending.Add(int64(1))

	return nil
}
