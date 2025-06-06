package contract

import "time"

// CacheStrategy 缓存策略接口
type CacheStrategy interface {
	// ShouldCache 是否应该缓存
	ShouldCache(key string, value interface{}) bool
	// GetTTL 获取TTL
	GetTTL(key string, value interface{}) time.Duration
	// GetCacheKey 获取缓存键
	GetCacheKey(prefix string, parts ...string) string
	// ShouldWarmUp 是否启用预热
	ShouldWarmUp() bool
}

// Serializer 序列化器接口
type Serializer interface {
	Serialize(v interface{}) ([]byte, error)
	Deserialize(data []byte, v interface{}) error
	ContentType() string
}
