package contracts

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) ([]byte, error) // Err 返回可兼容为“不存在”
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	ScanDeleteByPrefix(ctx context.Context, prefix string) error
	Incr(ctx context.Context, key string) (int64, error) // 原子自增（用于标签版本）
}
