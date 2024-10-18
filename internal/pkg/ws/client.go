package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dysodeng/app/internal/pkg/ws/message"

	"github.com/gorilla/websocket"
)

const (
	// 推送消息的超时时间
	writeWait = 10 * time.Second

	// 允许客户端发送的最大消息大小
	maxMessageSize = 512
)

// Client 是一个中间人，负责 WebSocket 连接和 Hub 之间的通信
type Client struct {
	// 底层的 WebSocket 连接
	conn *websocket.Conn

	// 缓冲发送消息的通道
	send chan message.WsMessage

	// 关联的用户 id
	userId string

	// 心跳计时器
	heartbeatTicker *time.Ticker
}

// 连接关闭时的处理函数
// 正常的断开不做处理，非正常的断开打印日志
func closeHandler(code int, text string) error {
	if code >= 1002 {
		log.Println("connection close: ", code, text)
	}
	return nil
}

// readMessage 从 WebSocket 连接中读取消息
//
// 该方法在一个独立的协程中运行，我们保证了每个连接只有一个 reader。
// 该方法会丢弃所有客户端传来的消息，如果需要接收可以在这里进行处理。
func (c *Client) readMessage() {
	defer func() {
		// unregister 为无缓冲通道，下面这一行会阻塞，
		// 直到 hub.run 中的 <-h.unregister 语句执行
		HubBus.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Time{}) // 永不超时
	for {
		// 从客户端接收消息
		messageType, body, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if HubBus.messageHandler != nil {
			err = HubBus.messageHandler.Handler(c.userId, messageType, body)
			if err != nil {
				log.Println("handler error: ", err)
			}
		}

	}
}

// writeMessage 负责推送消息给 WebSocket 客户端
//
// 该方法在一个独立的协程中运行，我们保证了每个连接只有一个 writer。
// Client 会从 send 请求中获取消息，然后在这个方法中推送给客户端。
func (c *Client) writeMessage() {
	defer func() {
		_ = c.conn.Close()
	}()

	// 从 send 通道中获取消息，然后推送给客户端
	for {
		messageItem, ok := <-c.send

		// 设置写超时时间
		_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		// c.send 这个通道已经被关闭了
		if !ok {
			HubBus.pending.Add(int64(-1 * len(c.send)))
			return
		}

		msg, _ := json.Marshal(messageItem.Message)
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			HubBus.errorHandler(messageItem, err)
			HubBus.pending.Add(int64(-1 * len(c.send)))
			return
		}

		HubBus.pending.Add(int64(-1))
	}
}

// heartbeat 发送心跳
func (c *Client) heartbeat() {
	for range c.heartbeatTicker.C {
		// 设置写超时时间
		_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

		if err := c.conn.WriteMessage(
			websocket.TextMessage,
			StringToBytes(fmt.Sprintf(`{"type":"%s"}`, message.TypeHeartbeat)),
		); err != nil {
			log.Println("heartbeat error: ", err)
			return
		}
	}
}
