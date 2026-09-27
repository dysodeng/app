package wx

import (
	"fmt"
	"sync"

	"github.com/goairix/wx/v2/miniapp"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/pkg/redis"
	wxCache "github.com/dysodeng/app/internal/infrastructure/pkg/wx/cache"
)

var (
	miniProgram       *miniapp.Client
	miniProgramConfig miniapp.Config
	miniProgramMu     sync.Mutex
)

// MiniProgram 返回微信小程序 SDK 客户端。
func MiniProgram() (*miniapp.Client, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		return nil, fmt.Errorf("wx miniapp: config is not loaded")
	}
	requested := miniapp.Config{
		AppID:     cfg.ThirdParty.Wx.MiniProgram.AppId,
		AppSecret: cfg.ThirdParty.Wx.MiniProgram.Secret,
	}

	miniProgramMu.Lock()
	defer miniProgramMu.Unlock()
	if miniProgram != nil && miniProgramConfig == requested {
		return miniProgram, nil
	}
	client, err := miniapp.NewClient(
		requested,
		miniapp.WithCache(wxCache.NewRedis(redis.CacheClient())),
	)
	if err != nil {
		return nil, err
	}
	miniProgram = client
	miniProgramConfig = requested
	return client, nil
}
