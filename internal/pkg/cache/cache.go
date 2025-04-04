package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/dysodeng/app/internal/config"
	"github.com/pkg/errors"
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
		switch config.Cache.Driver {
		case "memory": // 内存
			cache = NewMemoryCache()
		case "redis": // redis
			cache = NewRedis(
				fmt.Sprintf("%s:%s", config.Redis.Cache.Host, config.Redis.Cache.Port),
				config.Redis.Cache.Password,
				config.Redis.Cache.KeyPrefix,
				config.Redis.Cache.DB,
			)
		default:
			panic("缓存驱动错误")
		}
	})
	return cache, nil
}
