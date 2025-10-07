package driver

import (
	"context"
	"io"
)

// SpeechTextConverter 语音->文字转换器
type SpeechTextConverter interface {
	// ConvertSpeechToText 一次性转换语音到文字
	ConvertSpeechToText(ctx context.Context, reader io.Reader) (*Text, error)

	// StartStreamRecognition 开始流式语音识别
	// 返回一个通道用于接收识别结果和一个用于发送错误的通道
	StartStreamRecognition(ctx context.Context, config SpeechRecognitionConfig) (chan *Text, chan error, error)

	// SendAudioStream 发送音频流数据
	// 用于向已启动的流式识别会话发送音频数据
	SendAudioStream(ctx context.Context, data []byte) error

	// StopStreamRecognition 停止流式语音识别
	// 通知服务停止接收音频并完成最终识别
	StopStreamRecognition(ctx context.Context) error

	GetInitCompleteChan() chan struct{}
}

// TextSpeechConverter 文字->语音转换器
type TextSpeechConverter interface {
	// ConvertTextToSpeech 将文本转换为语音
	// 输入文本内容，返回语音文件信息
	ConvertTextToSpeech(ctx context.Context, text string) (*Speech, error)

	// ConvertTextToSpeechWithConfig 使用自定义配置将文本转换为语音
	// 可以指定语音参数如音色、语速等
	ConvertTextToSpeechWithConfig(ctx context.Context, text string, config SpeechSynthesisConfig) (*Speech, error)

	// ConvertTextToSpeechStream 将文本转换为语音并写入流
	// 适用于需要直接处理音频流而非文件的场景
	ConvertTextToSpeechStream(ctx context.Context, text string, writer io.Writer) error
}
