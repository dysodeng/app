package cache

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contracts"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/serializer"
)

// TypedCache 基于泛型的强类型缓存，支持标签失效与单航班防击穿
type TypedCache[T any] struct {
	ns         string
	cache      contracts.Cache
	serializer contracts.Serializer[T]
	sf         SFGroup[T]
	defaultTTL time.Duration
}

// NewTypedCache 创建类型缓存
func NewTypedCache[T any](namespace string, cache contracts.Cache) *TypedCache[T] {
	s := serializer.NewJSONSerializer[T]()
	if config.GlobalConfig.Cache.Serializer == "msgpack" {
		s = serializer.NewMsgpackSerializer[T]()
	}
	return &TypedCache[T]{
		ns:         strings.TrimSpace(namespace),
		cache:      cache,
		serializer: s,
	}
}

func (c *TypedCache[T]) WithSerializer(serializer contracts.Serializer[T]) *TypedCache[T] {
	c.serializer = serializer
	return c
}

func (c *TypedCache[T]) WithDefaultTTL(ttl time.Duration) *TypedCache[T] {
	c.defaultTTL = ttl
	return c
}

func (c *TypedCache[T]) versionKey(tag string) string {
	return fmt.Sprintf("__cv:%s:tag:%s", c.ns, tag)
}

func (c *TypedCache[T]) buildKey(base string, tags []string) (string, error) {
	// 固定顺序，避免同一集合不同顺序造成缓存击穿
	sort.Strings(tags)
	parts := []string{fmt.Sprintf("ns:%s", c.ns), fmt.Sprintf("k:%s", base)}
	for _, tag := range tags {
		// 获取标签版本，不存在视为0
		verBytes, _ := c.cache.Get(context.Background(), c.versionKey(tag))
		ver := "0"
		if len(verBytes) > 0 {
			ver = string(verBytes)
		}
		parts = append(parts, fmt.Sprintf("t:%s:%s", tag, ver))
	}
	return strings.Join(parts, "|"), nil
}

// Get 从缓存获取
func (c *TypedCache[T]) Get(ctx context.Context, base string, tags ...string) (T, bool, error) {
	key, err := c.buildKey(base, tags)
	if err != nil {
		var zero T
		return zero, false, err
	}
	b, err := c.cache.Get(ctx, key)
	if err != nil || len(b) == 0 {
		var zero T
		return zero, false, err
	}
	val, err := c.serializer.Decode(b)
	return val, err == nil, err
}

// Set 写入缓存
func (c *TypedCache[T]) Set(ctx context.Context, base string, val T, ttl time.Duration, tags ...string) error {
	key, err := c.buildKey(base, tags)
	if err != nil {
		return err
	}
	b, err := c.serializer.Encode(val)
	if err != nil {
		return err
	}
	if ttl <= 0 {
		ttl = c.defaultTTL
	}
	return c.cache.Set(ctx, key, b, ttl)
}

// Delete 删除缓存
func (c *TypedCache[T]) Delete(ctx context.Context, base string, tags ...string) error {
	key, err := c.buildKey(base, tags)
	if err != nil {
		return err
	}
	return c.cache.Delete(ctx, key)
}

// InvalidateTags 标签失效（版本 +1），O(1) 完成，无需扫描
func (c *TypedCache[T]) InvalidateTags(ctx context.Context, tags ...string) error {
	for _, tag := range tags {
		if _, err := c.cache.Incr(ctx, c.versionKey(tag)); err != nil {
			return err
		}
	}
	return nil
}

// BatchDeleteByPrefix 前缀删除（用于紧急清理，或老组件兼容）
func (c *TypedCache[T]) BatchDeleteByPrefix(ctx context.Context, prefix string) error {
	return c.cache.ScanDeleteByPrefix(ctx, prefix)
}

// GetOrLoad 旁路缓存 + 单航班防击穿
func (c *TypedCache[T]) GetOrLoad(ctx context.Context, base string, ttl time.Duration, loader func(context.Context) (T, error), tags ...string) (T, error) {
	var zero T
	key, err := c.buildKey(base, tags)
	if err != nil {
		return zero, err
	}

	// 先读缓存
	if raw, err := c.cache.Get(ctx, key); err == nil && len(raw) > 0 {
		return c.serializer.Decode(raw)
	}

	// singleflight 防击穿
	val, err := c.sf.Do(ctx, key, func() (T, error) {
		// 双检：避免并发间隙重复加载
		if raw, e := c.cache.Get(ctx, key); e == nil && len(raw) > 0 {
			return c.serializer.Decode(raw)
		}

		// 真正加载
		v, e := loader(ctx)
		if e != nil {
			var z T
			return z, e
		}

		// 回写缓存
		if enc, e := c.serializer.Encode(v); e == nil {
			if ttl <= 0 {
				ttl = c.defaultTTL
			}
			_ = c.cache.Set(ctx, key, enc, ttl)
		}
		return v, nil
	})
	if err != nil {
		return zero, err
	}

	return val, nil
}
