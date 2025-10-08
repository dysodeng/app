package websocket

import (
	"context"
	"encoding/json"

	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
	"github.com/dysodeng/app/internal/infrastructure/shared/websocket"
	"github.com/dysodeng/app/internal/interfaces/websocket/speech"
)

// MessageBody WebSocket消息体基础结构
type MessageBody struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"` // 使用json.RawMessage替代map[string]any
}

// textMessageHandler websocket文本消息接收处理器
type textMessageHandler struct {
	websocket.UnimplementedTextMessageHandler
}

func NewTextMessageHandler() websocket.TextMessageHandler {
	return &textMessageHandler{}
}

func (handler *textMessageHandler) Handler(ctx context.Context, clientId, userId string, data []byte) error {
	spanCtx, span := trace.Tracer().Start(ctx, "websocket.message.Handler")
	defer span.End()

	var messageBody MessageBody
	err := json.Unmarshal(data, &messageBody)
	if err != nil {
		logger.Error(spanCtx, "消息解析失败", logger.AddField("原消息", string(data)), logger.ErrorField(err))
		return err
	}

	if messageBody.Type == "speech" {
		// 实时语音转文字
		return speech.HandleSpeechMessage(spanCtx, clientId, userId, messageBody.Body)
	}

	return nil
}

// binaryMessageHandler websocket二进制消息接收处理器
type binaryMessageHandler struct {
	websocket.UnimplementedBinaryMessageHandler
}

func NewBinaryMessageHandler() websocket.BinaryMessageHandler {
	return &binaryMessageHandler{}
}

func (handler *binaryMessageHandler) Handler(ctx context.Context, clientId, userId string, data []byte) error {
	spanCtx, span := trace.Tracer().Start(ctx, "websocket.binaryMessageHandler.Handler")
	defer span.End()

	// 直接将二进制数据发送到语音识别会话
	err := speech.SessionManager.SendAudioData(ctx, clientId, data)
	if err != nil {
		logger.Error(spanCtx, "发送二进制音频数据失败", logger.ErrorField(err))
		return err
	}

	return nil
}
