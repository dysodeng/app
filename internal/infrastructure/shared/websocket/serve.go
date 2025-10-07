package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/resource"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
	"github.com/dysodeng/app/internal/infrastructure/shared/websocket/message"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Serve 处理 WebSocket 连接请求
func Serve(writer http.ResponseWriter, request *http.Request) {

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(writer, "upgrade error: %s", err.Error())
		return
	}

	// 认证失败的时候，返回错误信息，并断开连接
	userAuth, err := HubBus.authenticator.Authenticate(request)
	if err != nil || userAuth == nil {
		if err == nil {
			err = errors.New("用户信息为空")
		}
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
		_ = conn.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf(`{"type": "%s", "error": "authenticate error: %s"}`, message.TypeError, err.Error())),
		)
		_ = conn.Close()
		return
	}

	userId := userAuth["user_id"].(string)

	// 注册 Client
	clientId, err := uuid.NewV7()
	if err != nil {
		clientId = uuid.New()
	}
	client := &Client{
		conn:            conn,
		send:            make(chan message.WsMessage, bufferSize),
		userId:          userId,
		clientId:        clientId.String(),
		heartbeatTicker: time.NewTicker(time.Second * 5),
		tracer:          trace.TracerProvider().Tracer(resource.ServiceName() + ".websocket"),
	}
	client.conn.SetCloseHandler(closeHandler)

	// register 无缓冲，下面这一行会阻塞，直到 hub.run 中的 <-h.register 语句执行
	// 这样可以保证 register 成功之后才会启动读写协程
	HubBus.register <- client

	logger.Debug(request.Context(), "websocket.connection", logger.Field{Key: "client_id", Value: client.clientId}, logger.Field{Key: "user_id", Value: userId})

	// 启动读写协程
	go client.writeMessage()
	go client.readMessage()

	// 启动心跳协程
	go client.heartbeat()
}
