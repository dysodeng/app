package cache

import (
	"time"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contracts"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/driver"
)

// NewCacheDriver 创建缓存驱动
func NewCacheDriver(driverName string) contracts.Cache {
	switch driverName {
	case "redis":
		return driver.NewRedisCache()
	default:
		return driver.NewMemoryCache()
	}
}

// NewTypedCacheWith 优雅构建 TypedCache，包含默认 TTL 与序列化器选择
// driverName: "memory" 或 "redis"
// namespace: 缓存命名空间前缀
// ttl: 默认TTL（<=0 则不设置默认TTL）
func NewTypedCacheWith[T any](driverName, namespace string, ttl time.Duration) *TypedCache[T] {
	cacheDriver := NewCacheDriver(driverName)
	tc := NewTypedCache[T](namespace, cacheDriver)
	if ttl > 0 {
		tc.WithDefaultTTL(ttl)
	}
	return tc
}
