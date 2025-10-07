package driver

// Provider 供应商
type Provider string

const (
	ProviderAliyun Provider = "aliyun"
)

type Driver string

const (
	TTSDriver    Driver = "tts"
	SpeechDriver Driver = "speech"
)

// TTSConfig tts配置
type TTSConfig struct {
	TaskID string // 任务ID
	Voice  string // 语音音色
	Format string // 语音文件格式
}

// Text 语音转文字结果
type Text struct {
	Content string // 转换后的文本内容
	IsEnd   bool   // 是否是最终结果
}

// Speech 文字转语音结果
type Speech struct {
	AudioURL string `json:"audio"`    // 语音文件URL
	Duration int    `json:"duration"` // 语音时长(毫秒)
	Format   string `json:"format"`   // 音频格式
}

// SpeechRecognitionConfig 语音识别配置
type SpeechRecognitionConfig struct {
	Format                      string // 音频格式，如PCM、WAV、OPUS等
	SampleRate                  int    // 采样率，默认16000 Hz
	EnableIntermediateResult    bool   // 是否返回中间结果
	EnablePunctuationPrediction bool   // 是否开启标点预测
	EnableITN                   bool   // 是否开启数字转换(ITN)
}

// SpeechSynthesisConfig 语音合成配置
type SpeechSynthesisConfig struct {
	Voice      string // 发音人
	Format     string // 音频格式
	SampleRate int    // 采样率
	Volume     int    // 音量
	SpeechRate int    // 语速
	PitchRate  int    // 语调
}
