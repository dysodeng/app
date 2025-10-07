package cache

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
)

var (
	ErrKeyExpired  = errors.New("key expired")
	ErrKeyNotExist = errors.New("key not exist")
)

// Cache 缓存接口
type Cache interface {
	IsExist(key string) bool
	Get(key string) (string, error)
	Put(key string, value string, expiration time.Duration) error
	Delete(key string) error
	BatchDelete(prefix string) error
}

var cache Cache
var cacheInstanceOnce sync.Once

// NewCache 创建缓存实例
func NewCache() (Cache, error) {
	cacheInstanceOnce.Do(func() {
		switch config.GlobalConfig.Cache.Driver {
		case "memory": // 内存
			cache = NewMemoryCache()
		case "redis": // redis
			cache = NewRedisWithClient(redis.CacheClient(), "")
		default:
			panic("缓存驱动错误")
		}
	})
	return cache, nil
}
