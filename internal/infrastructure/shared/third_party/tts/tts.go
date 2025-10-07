package tts

import (
	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/shared/third_party/tts/driver"
	"github.com/dysodeng/app/internal/infrastructure/shared/third_party/tts/driver/aliyun"
)

// CreateTextToSpeech 创建文字转语音服务
func CreateTextToSpeech(provider driver.Provider, config driver.TTSConfig) (driver.TextSpeechConverter, error) {
	switch provider {
	case driver.ProviderAliyun:
		return aliyun.NewTTS(config)
	}
	return nil, errors.New("文字转语音服务提供商类型错误")
}

// CreateSpeechToText 创建语音转文字服务
func CreateSpeechToText(provider driver.Provider) (driver.SpeechTextConverter, error) {
	switch provider {
	case driver.ProviderAliyun:
		return aliyun.NewSpeech()
	}
	return nil, errors.New("语音转文字服务提供商类型错误")
}
