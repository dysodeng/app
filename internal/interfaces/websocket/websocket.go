package websocket

import "github.com/dysodeng/app/internal/infrastructure/shared/websocket"

// WebSocket 消息处理聚合器
type WebSocket struct {
	textMessageHandler   websocket.TextMessageHandler
	binaryMessageHandler websocket.BinaryMessageHandler
}

func NewWebSocket(
	textMessageHandler websocket.TextMessageHandler,
	binaryMessageHandler websocket.BinaryMessageHandler,
) *WebSocket {
	return &WebSocket{
		textMessageHandler:   textMessageHandler,
		binaryMessageHandler: binaryMessageHandler,
	}
}

func (ws *WebSocket) TextMessageHandler() websocket.TextMessageHandler {
	return ws.textMessageHandler
}

func (ws *WebSocket) BinaryMessageHandler() websocket.BinaryMessageHandler {
	return ws.binaryMessageHandler
}
