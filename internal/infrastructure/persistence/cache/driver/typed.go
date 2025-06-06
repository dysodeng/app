package driver

import (
	"context"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
)

type typedCache struct {
	cache      contract.Cache
	serializer contract.Serializer
}

func NewTypedCache(cache contract.Cache, serializer contract.Serializer) contract.TypedCache {
	return &typedCache{
		cache:      cache,
		serializer: serializer,
	}
}

func (tc *typedCache) GetObject(ctx context.Context, key string, obj interface{}) error {
	data, err := tc.cache.Get(ctx, key)
	if err != nil {
		return err
	}
	return tc.serializer.Deserialize([]byte(data), obj)
}

func (tc *typedCache) PutObject(ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	data, err := tc.serializer.Serialize(obj)
	if err != nil {
		return err
	}
	return tc.cache.Put(ctx, key, string(data), expiration)
}

func (tc *typedCache) DeleteObject(ctx context.Context, key string) error {
	return tc.cache.Delete(ctx, key)
}
