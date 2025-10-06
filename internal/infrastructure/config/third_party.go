package config

import (
	"github.com/spf13/viper"
)

// ThirdParty 第三方服务
type ThirdParty struct {
	Wx  wx  `mapstructure:"wx"`
	TTS tts `mapstructure:"tts"`
}

type wx struct {
	MiniProgram wxAccount `mapstructure:"mini_program"`
}

// wxAccount 微信公众账户
type wxAccount struct {
	AppId      string `mapstructure:"app_id"`
	Secret     string `mapstructure:"secret"`
	OriginalID string `mapstructure:"original_id"`
	EnvVersion string `mapstructure:"env_version"`
}

// tts 语音服务
type tts struct {
	TTSProvider       string   `mapstructure:"tts_provider"`
	TTSVoice          string   `mapstructure:"tts_voice"`
	TTSVoiceFormat    string   `mapstructure:"tts_voice_format"`
	SpeechProvider    string   `mapstructure:"speech_provider"`
	SpeechVoiceFormat string   `mapstructure:"speech_voice_format"`
	Provider          provider `mapstructure:"provider"`
}

type provider struct {
	Aliyun ttsProvider `mapstructure:"aliyun"`
}

type ttsProvider struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	AppKey          string `mapstructure:"app_key"`
}

func thirdPartyBindEnv(d *viper.Viper) {
	_ = d.BindEnv("wx.mini_program.app_id", "WX_MINI_PROGRAM_APP_ID")
	_ = d.BindEnv("wx.mini_program.secret", "WX_MINI_PROGRAM_SECRET")
	_ = d.BindEnv("wx.mini_program.original_id", "WX_MINI_PROGRAM_ORIGINAL_ID")
	_ = d.BindEnv("wx.mini_program.env_version", "WX_MINI_PROGRAM_ENV_VERSION")
	_ = d.BindEnv("tts.tts_provider", "TTS_PROVIDER")
	_ = d.BindEnv("tts.speech_provider", "TTS_SPEECH_PROVIDER")
	_ = d.BindEnv("tts.provider.aliyun.access_key_id", "TTS_ALIYUN_ISI_ACCESS_KEY_ID")
	_ = d.BindEnv("tts.provider.aliyun.access_key_secret", "TTS_ALIYUN_ISI_ACCESS_KEY_SECRET")
	_ = d.BindEnv("tts.provider.aliyun.app_key", "TTS_ALIYUN_ISI_APP_KEY")
}
