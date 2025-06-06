package strategy

import (
	"strings"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
)

type userStrategy struct {
	defaultTTL time.Duration
}

func NewUserStrategy() contract.CacheStrategy {
	return &userStrategy{
		defaultTTL: 1 * time.Hour,
	}
}

func (s *userStrategy) ShouldCache(key string, value interface{}) bool {
	// 用户信息都应该缓存
	return true
}

func (s *userStrategy) GetTTL(key string, value interface{}) time.Duration {
	if strings.Contains(key, "info") {
		// 用户信息缓存1小时
		return 1 * time.Hour
	}
	if strings.Contains(key, "list") {
		// 用户列表缓存30分钟
		return 30 * time.Minute
	}
	return s.defaultTTL
}

func (s *userStrategy) GetCacheKey(prefix string, parts ...string) string {
	allParts := append([]string{prefix}, parts...)
	return strings.Join(allParts, ":")
}

func (s *userStrategy) ShouldWarmUp() bool {
	return false
}
