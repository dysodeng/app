package ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dysodeng/app/internal/pkg/ws/message"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebsocketServe 处理 WebSocket 连接请求
func WebsocketServe(writer http.ResponseWriter, request *http.Request) {

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(fmt.Sprintf("upgrade error: %s", err.Error())))
		return
	}

	// 认证失败的时候，返回错误信息，并断开连接
	userId, err := HubBus.authenticator.Authenticate(request)
	if err != nil {
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
		_ = conn.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf(`{"type": "%s", "error": "authenticate error: %s"}`, message.TypeError, err.Error())),
		)
		_ = conn.Close()
		return
	}

	// 注册 Client
	client := &Client{
		conn:            conn,
		send:            make(chan message.WsMessage, bufferSize),
		userId:          userId,
		heartbeatTicker: time.NewTicker(time.Second * 5),
	}
	client.conn.SetCloseHandler(closeHandler)

	// register 无缓冲，下面这一行会阻塞，直到 hub.run 中的 <-h.register 语句执行
	// 这样可以保证 register 成功之后才会启动读写协程
	HubBus.register <- client

	// 启动读写协程
	go client.writeMessage()
	go client.readMessage()

	// 启动心跳协程
	go client.heartbeat()
}
