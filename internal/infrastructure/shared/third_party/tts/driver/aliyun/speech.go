package aliyun

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
	"time"

	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/helper"
	"github.com/dysodeng/app/internal/infrastructure/shared/third_party/tts/driver"
)

type SpeechResult struct {
	Header struct {
		Namespace  string `json:"namespace"`
		Name       string `json:"name"`
		Status     int    `json:"status"`
		MessageID  string `json:"message_id"`
		TaskID     string `json:"task_id"`
		StatusText string `json:"status_text"`
	} `json:"header"`
	Payload struct {
		Result   string `json:"result"`
		Duration int    `json:"duration"`
	} `json:"payload"`
}

// Speech 阿里云语音识别实现
type Speech struct {
	recognition      *nls.SpeechRecognition
	resultChannel    chan *driver.Text
	errorChannel     chan error
	connectionToken  string
	appKey           string
	logger           *nls.NlsLogger
	started          bool
	mutex            sync.Mutex
	initCompleteChan chan struct{}
}

// NewSpeech 创建阿里云语音识别实例
func NewSpeech() (*Speech, error) {
	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	logger := nls.NewNlsLogger(os.Stderr, "SpeechRecognition: ", log.LstdFlags|log.Lmicroseconds)

	return &Speech{
		resultChannel:   make(chan *driver.Text),
		errorChannel:    make(chan error),
		connectionToken: token.Token.Id,
		appKey:          config.GlobalConfig.ThirdParty.TTS.Provider.Aliyun.AppKey,
		logger:          logger,
		started:         false,
	}, nil
}

// ConvertSpeechToText 实现一次性语音转文字
func (s *Speech) ConvertSpeechToText(ctx context.Context, reader io.Reader) (*driver.Text, error) {
	// 读取所有音频数据
	audioData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// 创建上下文，设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 启动流式识别
	textChan, errChan, err := s.StartStreamRecognition(ctx, driver.SpeechRecognitionConfig{
		Format:                      config.GlobalConfig.ThirdParty.TTS.SpeechVoiceFormat,
		SampleRate:                  16000,
		EnableIntermediateResult:    false,
		EnablePunctuationPrediction: true,
		EnableITN:                   true,
	})
	if err != nil {
		return nil, err
	}

	// 确保资源清理
	defer func() {
		_ = s.StopStreamRecognition(ctx)
	}()

	// 发送音频数据
	err = s.SendAudioStream(ctx, audioData)
	if err != nil {
		return nil, err
	}

	// 停止识别
	err = s.StopStreamRecognition(ctx)
	if err != nil {
		return nil, err
	}

	// 等待结果
	var finalResult *driver.Text

	// 使用单独的通道接收最终结果
	resultChan := make(chan *driver.Text, 1)
	locErrChan := make(chan error, 1)

	go func() {
		for {
			select {
			case result, ok := <-textChan:
				if !ok {
					// 通道已关闭
					if finalResult != nil {
						resultChan <- finalResult
					}
					return
				}
				if result.IsEnd {
					resultChan <- result
					return
				}
				finalResult = result
			case err, ok := <-errChan:
				if !ok {
					return
				}
				locErrChan <- err
				return
			case <-ctx.Done():
				locErrChan <- errors.New("识别超时")
				return
			}
		}
	}()

	// 等待最终结果或错误
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-locErrChan:
		return nil, err
	case <-ctx.Done():
		return nil, errors.New("识别超时")
	}
}

// StartStreamRecognition 开始流式语音识别
func (s *Speech) StartStreamRecognition(ctx context.Context, config driver.SpeechRecognitionConfig) (chan *driver.Text, chan error, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.started {
		return nil, nil, errors.New("语音识别已经启动")
	}

	// 创建连接配置
	connConfig := nls.NewConnectionConfigWithToken(nls.DEFAULT_URL, s.appKey, s.connectionToken)

	// 创建结果和错误通道
	s.resultChannel = make(chan *driver.Text, 10)
	s.errorChannel = make(chan error, 1)
	s.initCompleteChan = make(chan struct{}, 1)

	// 创建语音识别实例
	recognition, err := nls.NewSpeechRecognition(
		connConfig,
		s.logger,
		s.onTaskFailed,
		s.onStarted,
		s.onResultChanged,
		s.onCompleted,
		s.onClose,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	s.recognition = recognition

	// 设置识别参数
	param := nls.DefaultSpeechRecognitionParam()
	if config.Format != "" {
		param.Format = config.Format
	}
	if config.SampleRate > 0 {
		param.SampleRate = config.SampleRate
	}
	param.EnableIntermediateResult = config.EnableIntermediateResult
	param.EnablePunctuationPrediction = config.EnablePunctuationPrediction
	param.EnableInverseTextNormalization = config.EnableITN

	// 启动识别
	_, err = s.recognition.Start(param, nil)
	if err != nil {
		return nil, nil, err
	}

	s.started = true

	return s.resultChannel, s.errorChannel, nil
}

// SendAudioStream 发送音频流数据
func (s *Speech) SendAudioStream(ctx context.Context, data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.started || s.recognition == nil {
		return errors.New("语音识别未启动")
	}

	// 发送音频数据，最多重试3次
	var err error
	for i := 0; i < 3; i++ {
		err = s.recognition.SendAudioData(data)
		if err == nil {
			return nil
		}

		// 如果是连接关闭错误，不再重试
		if err.Error() == "use of closed network connection" {
			break
		}

		// 短暂等待后重试
		time.Sleep(10 * time.Millisecond)
	}

	return err
}

// StopStreamRecognition 停止流式语音识别
func (s *Speech) StopStreamRecognition(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.started || s.recognition == nil {
		return nil // 已经停止，不返回错误
	}

	// 停止发送音频
	_, err := s.recognition.Stop()
	// 即使出错也继续清理资源

	// 关闭连接
	s.recognition.Shutdown()
	s.started = false

	// 主动发送一个最终结果，确保客户端收到is_end=true的消息
	if s.resultChannel != nil {
		select {
		case s.resultChannel <- &driver.Text{
			Content: "", // 可以为空或者使用最后一次的识别结果
			IsEnd:   true,
		}:
			// 成功发送结果
		default:
			// 通道已满或已关闭，记录日志
			s.logger.Printf("Failed to send final result to channel after stop")
		}
	}

	return err
}

func (s *Speech) GetInitCompleteChan() chan struct{} {
	return s.initCompleteChan
}

// 回调函数

func (s *Speech) onTaskFailed(text string, param interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.errorChannel != nil {
		select {
		case s.errorChannel <- errors.New(text):
			// 成功发送错误
		default:
			// 通道已满或已关闭，记录日志
			s.logger.Printf("Failed to send error to channel: %s", text)
		}
	}
}

func (s *Speech) onStarted(text string, param interface{}) {
	// 识别已启动
	s.logger.Printf("Recognition started: %s", text)

	if s.initCompleteChan != nil {
		select {
		case s.initCompleteChan <- struct{}{}: // 成功发送信号
		default: // 通道已满或已关闭
		}
	}
}

func (s *Speech) onResultChanged(text string, param interface{}) {
	// 中间结果
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var result SpeechResult
	if err := sonic.Unmarshal(helper.StringToBytes(text), &result); err != nil {
		s.logger.Printf("Failed to parse result JSON: %v", err)
		return
	}

	if s.resultChannel != nil {
		select {
		case s.resultChannel <- &driver.Text{
			Content: result.Payload.Result,
			IsEnd:   false,
		}:
			// 成功发送结果
		default:
			// 通道已满或已关闭，记录日志
			s.logger.Printf("Failed to send intermediate result to channel")
		}
	}
}

func (s *Speech) onCompleted(text string, param interface{}) {
	// 最终结果
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var result SpeechResult
	if err := sonic.Unmarshal(helper.StringToBytes(text), &result); err != nil {
		s.logger.Printf("Failed to parse result JSON: %v", err)
		return
	}

	if s.resultChannel != nil {
		select {
		case s.resultChannel <- &driver.Text{
			Content: result.Payload.Result,
			IsEnd:   true,
		}:
			// 成功发送结果
		default:
			// 通道已满或已关闭，记录日志
			s.logger.Printf("Failed to send final result to channel")
		}
	}
}

func (s *Speech) onClose(param interface{}) {
	// 连接关闭，关闭通道
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 安全关闭通道
	if s.resultChannel != nil {
		close(s.resultChannel)
		s.resultChannel = nil
	}

	if s.errorChannel != nil {
		close(s.errorChannel)
		s.errorChannel = nil
	}

	s.logger.Printf("Connection closed")
}
