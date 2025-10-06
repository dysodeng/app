package wx

import (
	"sync"

	"github.com/dysodeng/wx/mini_program"
	"github.com/dysodeng/wx/support/cache"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
	wxCache "github.com/dysodeng/app/internal/infrastructure/shared/wx/cache"
)

var (
	cacheItem       cache.Cache
	miniProgram     *mini_program.MiniProgram
	miniProgramOnce sync.Once
)

// initMiniProgram 初始化微信小程序sdk
func initMiniProgram() *mini_program.MiniProgram {
	miniProgramOnce.Do(func() {
		cacheItem = wxCache.NewRedis(redis.CacheClient())
		miniProgram = mini_program.New(
			config.GlobalConfig.ThirdParty.Wx.MiniProgram.AppId,
			config.GlobalConfig.ThirdParty.Wx.MiniProgram.Secret,
			"",
			"",
			mini_program.WithCache(cacheItem),
		)
	})
	return miniProgram
}

func MiniProgram() *mini_program.MiniProgram {
	return initMiniProgram()
}
