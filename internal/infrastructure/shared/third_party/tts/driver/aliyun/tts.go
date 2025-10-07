package aliyun

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"time"

	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
	"github.com/dysodeng/app/internal/infrastructure/shared/third_party/tts/driver"
	"github.com/dysodeng/app/internal/infrastructure/shared/third_party/tts/driver/helper"
)

// TTS 文字转语音
type TTS struct {
	speechSynthesis *nls.SpeechSynthesis
	Reader          io.Reader
	config          driver.TTSConfig
}

type userParam struct {
	Writer io.Writer
	Logger *nls.NlsLogger
	TaskID string
}

func NewTTS(cfg driver.TTSConfig) (*TTS, error) {
	if cfg.TaskID == "" {
		return nil, errors.New("缺少TaskID")
	}

	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	ttsConfig := nls.NewConnectionConfigWithToken(
		nls.DEFAULT_URL,
		config.GlobalConfig.ThirdParty.TTS.Provider.Aliyun.AppKey,
		token.Token.Id,
	)

	logger := nls.NewNlsLogger(os.Stderr, "TaskID-"+cfg.TaskID+": ", log.LstdFlags|log.Lmicroseconds)
	logger.SetLogSil(false)
	if config.GlobalConfig.App.Debug {
		logger.SetDebug(true)
	}

	fileBytes := bytes.NewBuffer([]byte{})

	ttsUserParam := new(userParam)
	ttsUserParam.TaskID = cfg.TaskID
	ttsUserParam.Logger = logger
	ttsUserParam.Writer = fileBytes

	tts := &TTS{
		config: cfg,
		Reader: fileBytes,
	}

	speechSynthesis, err := nls.NewSpeechSynthesis(ttsConfig, logger, false,
		tts.onTaskFailed, tts.onSynthesisResult, nil,
		tts.onCompleted, tts.onClose, ttsUserParam)
	if err != nil {
		return nil, err
	}

	tts.speechSynthesis = speechSynthesis

	return tts, nil
}

// Start 开始任务
func (tts *TTS) Start(text string) (chan bool, error) {
	if text == "" {
		return nil, errors.New("文本为空")
	}
	param := nls.DefaultSpeechSynthesisParam()
	if tts.config.Voice != "" {
		param.Voice = tts.config.Voice
	}
	if tts.config.Format != "" {
		param.Format = tts.config.Format
	}

	ch, err := tts.speechSynthesis.Start(text, param, nil)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

// Wait 等待任务完成
func (tts *TTS) Wait(ch chan bool, onCompleted func(reader io.Reader), onFailed func(err error)) {
	select {
	case done := <-ch:
		if !done {
			tts.speechSynthesis.Shutdown()
			onFailed(errors.New("Wait Failed"))
			return
		}
		tts.speechSynthesis.Shutdown()
		onCompleted(tts.Reader)
	case <-time.After(60 * time.Second):
		tts.speechSynthesis.Shutdown()
		onFailed(errors.New("Wait timeout"))
		return
	}
}

func (tts *TTS) params(param interface{}) (*userParam, error) {
	p, ok := param.(*userParam)
	if !ok {
		return nil, errors.New("invalid params")
	}
	return p, nil
}

// onTaskFailed 任务失败事件回调
func (tts *TTS) onTaskFailed(text string, param interface{}) {
	// p, err := tts.params(param)
	// if err != nil {
	//	log.Println(err)
	//	return
	// }

	// p.Logger.Printf("TaskFailed: %s", text)
}

// onSynthesisResult 音频流回调
func (tts *TTS) onSynthesisResult(data []byte, param interface{}) {
	p, err := tts.params(param)
	if err != nil {
		log.Println(err)
		return
	}
	_, _ = p.Writer.Write(data)
}

// onCompleted 完成事件回调
func (tts *TTS) onCompleted(text string, param interface{}) {
	// p, err := tts.params(param)
	// if err != nil {
	//	log.Println(err)
	//	return
	// }

	// p.Logger.Println("onCompleted:", text)
}

// onClose 关闭链接事件回调
func (tts *TTS) onClose(param interface{}) {
	// p, err := tts.params(param)
	// if err != nil {
	//	log.Println(err)
	//	return
	// }

	// p.Logger.Println("onClosed")
}

// 实现 TextSpeechConverter 接口

// ConvertTextToSpeech 将文本转换为语音
func (tts *TTS) ConvertTextToSpeech(ctx context.Context, text string) (*driver.Speech, error) {
	ch, err := tts.Start(text)
	if err != nil {
		return nil, err
	}

	speech := &driver.Speech{}
	var convertErr error

	uploader := storage.Instance().FileSystem().Uploader()

	tts.Wait(ch, func(reader io.Reader) {
		// 上传文件
		path, _ := helper.GenerateFilePath(".wav")
		_ = uploader.Upload(ctx, path, reader)

		speech.Format = tts.config.Format
		speech.AudioURL = storage.Instance().FullUrl(ctx, path)
		speech.Duration = 0
	}, func(err error) {
		convertErr = err
	})

	if convertErr != nil {
		return nil, convertErr
	}

	return speech, nil
}

// ConvertTextToSpeechWithConfig 使用自定义配置将文本转换为语音
func (tts *TTS) ConvertTextToSpeechWithConfig(ctx context.Context, text string, config driver.SpeechSynthesisConfig) (*driver.Speech, error) {
	// 保存原始配置
	originalVoice := tts.config.Voice
	originalFormat := tts.config.Format

	// 应用新配置
	if config.Voice != "" {
		tts.config.Voice = config.Voice
	}
	if config.Format != "" {
		tts.config.Format = config.Format
	}

	// 转换
	speech, err := tts.ConvertTextToSpeech(ctx, text)

	// 恢复原始配置
	tts.config.Voice = originalVoice
	tts.config.Format = originalFormat

	return speech, err
}

// ConvertTextToSpeechStream 将文本转换为语音并写入流
func (tts *TTS) ConvertTextToSpeechStream(ctx context.Context, text string, writer io.Writer) error {
	ch, err := tts.Start(text)
	if err != nil {
		return err
	}

	var convertErr error

	tts.Wait(ch, func(reader io.Reader) {
		// 将音频数据从reader复制到writer
		_, err = io.Copy(writer, reader)
		if err != nil {
			convertErr = err
		}
	}, func(err error) {
		convertErr = err
	})

	return convertErr
}
