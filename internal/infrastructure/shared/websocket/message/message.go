package message

// Type 消息类型
type Type string

const (
	TypeHeartbeat Type = "heartbeat" // 心跳消息
	TypeMessage   Type = "message"   // 业务消息
	TypeError     Type = "error"     // 错误消息
)

type Message struct {
	Type Type   `json:"type"`
	Data string `json:"data"`
}

type WsMessage struct {
	ClientID string  `json:"client_id"`
	Message  Message `json:"message"`
}
