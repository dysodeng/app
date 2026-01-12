package websocket

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/websocket/message"
)

const (
	// 推送消息的超时时间
	writeWait = 10 * time.Second

	// 允许客户端发送的最大消息大小
	maxMessageSize = 1024 * 1024
)

// Client 是一个中间人，负责 WebSocket 连接和 Hub 之间的通信
type Client struct {
	// 底层的 WebSocket 连接
	conn *websocket.Conn
	// 缓冲发送消息的通道
	send chan message.WsMessage
	// 关联的用户id
	userId string
	// 客户端id
	clientId string
	// 心跳计时器
	heartbeatTicker *time.Ticker
	// 链路追踪
	tracer trace.Tracer
}

// 连接关闭时的处理函数
// 正常的断开不做处理，非正常的断开打印日志
func closeHandler(code int, text string) error {
	logger.Debug(context.Background(), "websocket client close", logger.Field{Key: "code", Value: code}, logger.Field{Key: "text", Value: text})
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
		ctx := context.Background()
		spanCtx, span := c.tracer.Start(ctx, "websocket.ReadMessage")

		// 从客户端接收消息
		messageType, body, err := c.conn.ReadMessage()
		if err != nil {
			if errorsTotal != nil {
				errorsTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("type", "read_message")))
			}
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			span.End()
			break
		}
		if messagesReceivedTotal != nil {
			messagesReceivedTotal.Add(ctx, 1, metric.WithAttributes(attribute.Int("message_type", messageType)))
		}
		span.SetAttributes(
			attribute.Int("ws.message.type", messageType),
			attribute.Int("ws.message.size", len(body)),
		)

		if messageType == websocket.BinaryMessage {
			// 二进制消息
			if HubBus.binaryMessageHandler != nil {
				err = HubBus.binaryMessageHandler.Handler(spanCtx, c.clientId, c.userId, body)
				if err != nil {
					if errorsTotal != nil {
						errorsTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("type", "handler_binary")))
					}
					span.SetStatus(codes.Error, err.Error())
					span.RecordError(err)
					log.Println("handler error: ", err)
				}
			}
		} else {
			// 文本消息
			if HubBus.textMessageHandler != nil {
				err = HubBus.textMessageHandler.Handler(spanCtx, c.clientId, c.userId, body)
				if err != nil {
					if errorsTotal != nil {
						errorsTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("type", "handler_text")))
					}
					span.SetStatus(codes.Error, err.Error())
					span.RecordError(err)
					log.Println("handler error: ", err)
				}
			}
		}
		span.End()
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

		msg, _ := sonic.Marshal(messageItem.Message)
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			if errorsTotal != nil {
				errorsTotal.Add(context.Background(), 1, metric.WithAttributes(attribute.String("type", "write_message")))
			}
			HubBus.errorHandler(messageItem, err)
			HubBus.pending.Add(int64(-1 * len(c.send)))
			return
		}

		if messagesSentTotal != nil {
			messagesSentTotal.Add(context.Background(), 1)
		}

		HubBus.pending.Add(int64(-1))
	}
}

// heartbeat 发送心跳
func (c *Client) heartbeat() {
	defer c.heartbeatTicker.Stop()

	for range c.heartbeatTicker.C {
		// 设置写超时时间
		_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

		if err := c.conn.WriteMessage(
			websocket.TextMessage,
			StringToBytes(fmt.Sprintf(`{"type":"%s"}`, message.TypeHeartbeat)),
		); err != nil {
			if errorsTotal != nil {
				errorsTotal.Add(context.Background(), 1, metric.WithAttributes(attribute.String("type", "heartbeat")))
			}
			// 记录详细的网络错误信息
			logger.Error(context.Background(), "websocket heartbeat failed",
				logger.Field{Key: "client_id", Value: c.clientId},
				logger.Field{Key: "user_id", Value: c.userId},
				logger.Field{Key: "error", Value: err.Error()},
			)

			// 网络错误时优雅关闭连接
			HubBus.unregister <- c
			return
		}
	}
}
