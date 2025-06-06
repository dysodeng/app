package contract

import (
	"context"
	"time"
)

// Cache 基础缓存接口
type Cache interface {
	// 基础操作
	IsExist(ctx context.Context, key string) bool
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	BatchDelete(ctx context.Context, prefix string) error

	// 高级操作
	GetWithTTL(ctx context.Context, key string) (string, time.Duration, error)
	Increment(ctx context.Context, key string, delta int64) (int64, error)
	Decrement(ctx context.Context, key string, delta int64) (int64, error)
	SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	FlushAll(ctx context.Context) error

	// 生命周期
	Close() error
}

// TypedCache 类型化缓存接口
type TypedCache interface {
	GetObject(ctx context.Context, key string, obj interface{}) error
	PutObject(ctx context.Context, key string, obj interface{}, expiration time.Duration) error
	DeleteObject(ctx context.Context, key string) error
}
