package cache

import (
	"fmt"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/driver"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/serializer"
	"github.com/dysodeng/app/internal/pkg/redis"
)

// Factory 缓存工厂
type Factory struct {
	cache      contract.Cache
	typedCache contract.TypedCache
}

func NewCacheFactory() (*Factory, error) {
	var cache contract.Cache

	switch config.Cache.Driver {
	case "redis":
		// 使用专门的缓存Redis客户端
		redisClient := redis.CacheClient()
		cache = driver.NewRedisCache(redisClient, config.Redis.Cache.KeyPrefix)
	case "memory":
		cache = driver.NewMemoryCache()
	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", config.Cache.Driver)
	}

	jsonSerializer := serializer.NewJSONSerializer()
	typedCache := driver.NewTypedCache(cache, jsonSerializer)

	return &Factory{
		cache:      cache,
		typedCache: typedCache,
	}, nil
}

func (f *Factory) GetCache() contract.Cache {
	return f.cache
}

func (f *Factory) GetTypedCache() contract.TypedCache {
	return f.typedCache
}
