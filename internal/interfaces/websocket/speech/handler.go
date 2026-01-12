package speech

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bytedance/sonic"

	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
	"github.com/dysodeng/app/internal/infrastructure/shared/websocket"
	"github.com/dysodeng/app/internal/infrastructure/shared/websocket/message"
)

// SessionManager 全局会话管理器
var SessionManager *Manager

// 初始化会话管理器
func init() {
	SessionManager = NewManager()
}

// SpeechMessage WebSocket语音消息
type SpeechMessage struct {
	Action string `json:"action"` // 动作：start, data, stop
	Data   string `json:"data"`   // Base64编码的音频数据
}

// SpeechResult 语音识别结果
type SpeechResult struct {
	Content string `json:"content"` // 识别文本
	IsEnd   bool   `json:"is_end"`  // 是否是最终结果
}

// HandleSpeechMessage 处理语音消息
func HandleSpeechMessage(ctx context.Context, clientID, userID string, body json.RawMessage) error {
	spanCtx, span := trace.Tracer().Start(ctx, "websocket.speech.HandleSpeechMessage")
	defer span.End()

	// 直接反序列化为SpeechMessage，避免中间的序列化步骤
	var speechMessage SpeechMessage
	if err := sonic.Unmarshal(body, &speechMessage); err != nil {
		logger.Error(spanCtx, "语音消息解析失败", logger.ErrorField(err))
		return err
	}

	switch speechMessage.Action {
	case "start":
		// 开始语音识别
		return handleStartAction(spanCtx, clientID, userID)
	case "stop":
		// 停止语音识别
		return handleStopAction(spanCtx, clientID)
	default:
		logger.Error(spanCtx, "未知的语音消息动作", logger.AddField("action", speechMessage.Action))
		return fmt.Errorf("未知的语音消息动作: %s", speechMessage.Action)
	}
}

// handleStartAction 处理开始语音识别动作
func handleStartAction(ctx context.Context, clientID, userID string) error {
	// 创建语音识别会话
	session, err := SessionManager.CreateSession(ctx, clientID, userID)
	if err != nil {
		logger.Error(ctx, "创建语音识别会话失败", logger.ErrorField(err))
		return err
	}

	// 启动goroutine监听识别结果
	go listenRecognitionResults(ctx, session)

	// 等待初始化完成或超时
	select {
	case <-session.InitCompleteChan:
		// 初始化完成
	case <-time.After(5 * time.Second):
		// 超时处理
		logger.Warn(ctx, "等待语音识别初始化完成超时")
	}

	// 发送开始识别成功消息
	result := SpeechResult{
		Content: "",
		IsEnd:   false,
	}
	resultBytes, _ := sonic.Marshal(result)
	return websocket.SendMessage(clientID, message.TypeMessage, string(resultBytes))
}

// handleStopAction 处理停止语音识别动作
func handleStopAction(ctx context.Context, clientID string) error {
	// 关闭会话
	if err := SessionManager.CloseSession(ctx, clientID); err != nil {
		logger.Error(ctx, "关闭语音识别会话失败", logger.ErrorField(err))
		return err
	}

	return nil
}

// listenRecognitionResults 监听识别结果
func listenRecognitionResults(ctx context.Context, session *Session) {
	for {
		select {
		case result, ok := <-session.ResultChannel:
			if !ok {
				// 通道已关闭
				return
			}

			// 发送识别结果给客户端
			speechResult := SpeechResult{
				Content: result.Content,
				IsEnd:   result.IsEnd,
			}
			resultBytes, _ := sonic.Marshal(speechResult)
			if err := websocket.SendMessage(session.ClientID, message.TypeMessage, string(resultBytes)); err != nil {
				logger.Error(ctx, "发送识别结果失败", logger.ErrorField(err))
			}

			// 如果是最终结果，关闭会话
			if result.IsEnd {
				_ = SessionManager.CloseSession(ctx, session.ClientID)
				return
			}

		case err, ok := <-session.ErrorChannel:
			if !ok {
				// 通道已关闭
				return
			}

			// 发送错误消息给客户端
			errMsg := fmt.Sprintf("语音识别错误: %s", err.Error())
			if err := websocket.SendMessage(session.ClientID, message.TypeError, errMsg); err != nil {
				logger.Error(ctx, "发送错误消息失败", logger.ErrorField(err))
			}

			// 关闭会话
			_ = SessionManager.CloseSession(ctx, session.ClientID)
			return
		}
	}
}
