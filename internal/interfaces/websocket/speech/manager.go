package speech

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/third_party/tts"
	"github.com/dysodeng/app/internal/infrastructure/shared/third_party/tts/driver"
)

// Session 语音识别会话
type Session struct {
	ClientID         string                     // 客户端ID
	UserID           string                     // 用户ID
	Converter        driver.SpeechTextConverter // 语音转文字服务
	ResultChannel    chan *driver.Text          // 结果通道
	ErrorChannel     chan error                 // 错误通道
	InitCompleteChan chan struct{}              // 初始化完成通道
	Active           bool                       // 会话是否活跃
}

// Manager 语音识别会话管理器
type Manager struct {
	sessions map[string]*Session // 客户端ID -> 会话
	mutex    sync.RWMutex        // 读写锁
}

// NewManager 创建语音识别会话管理器
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

// CreateSession 创建语音识别会话
func (m *Manager) CreateSession(ctx context.Context, clientID, userID string) (*Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查是否已存在会话
	if session, exists := m.sessions[clientID]; exists && session.Active {
		return session, nil
	}

	// 创建语音转文字服务
	converter, err := tts.CreateSpeechToText(driver.Provider(config.GlobalConfig.ThirdParty.TTS.SpeechProvider))
	if err != nil {
		return nil, err
	}

	// 启动流式识别
	resultChan, errChan, err := converter.StartStreamRecognition(ctx, driver.SpeechRecognitionConfig{
		Format:                      config.GlobalConfig.ThirdParty.TTS.SpeechVoiceFormat,
		SampleRate:                  16000,
		EnableIntermediateResult:    true,
		EnablePunctuationPrediction: true,
		EnableITN:                   true,
	})
	if err != nil {
		return nil, err
	}

	// 创建会话
	session := &Session{
		ClientID:         clientID,
		UserID:           userID,
		Converter:        converter,
		ResultChannel:    resultChan,
		ErrorChannel:     errChan,
		InitCompleteChan: converter.GetInitCompleteChan(),
		Active:           true,
	}

	// 保存会话
	m.sessions[clientID] = session

	return session, nil
}

// GetSession 获取语音识别会话
func (m *Manager) GetSession(ctx context.Context, clientID string) (*Session, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[clientID]
	if !exists || !session.Active {
		return nil, errors.New("语音识别会话不存在")
	}

	return session, nil
}

// CloseSession 关闭语音识别会话
func (m *Manager) CloseSession(ctx context.Context, clientID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.sessions[clientID]
	if !exists || !session.Active {
		return nil
	}

	// 停止流式识别
	err := session.Converter.StopStreamRecognition(ctx)

	// 标记会话为非活跃
	session.Active = false

	// 删除会话
	delete(m.sessions, clientID)

	return err
}

// SendAudioData 发送音频数据
func (m *Manager) SendAudioData(ctx context.Context, clientID string, data []byte) error {
	session, err := m.GetSession(ctx, clientID)
	if err != nil {
		return err
	}

	return session.Converter.SendAudioStream(ctx, data)
}
