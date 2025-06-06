package driver

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
)

type cacheItem struct {
	value      string
	expiration time.Time
	hasExpiry  bool
}

type memoryCache struct {
	mu        sync.RWMutex
	items     map[string]*cacheItem
	keyPrefix string
	cleaner   *time.Ticker
	stopChan  chan bool
}

func NewMemoryCache() contract.Cache {
	c := &memoryCache{
		items:    make(map[string]*cacheItem),
		cleaner:  time.NewTicker(time.Minute), // 每分钟清理一次过期项
		stopChan: make(chan bool),
	}

	// 启动清理协程
	go c.startCleaner()

	return c
}

func NewMemoryCacheWithPrefix(keyPrefix string) contract.Cache {
	c := NewMemoryCache().(*memoryCache)
	c.keyPrefix = keyPrefix
	return c
}

func (m *memoryCache) key(key string) string {
	if m.keyPrefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", m.keyPrefix, key)
}

func (m *memoryCache) startCleaner() {
	for {
		select {
		case <-m.cleaner.C:
			m.cleanExpired()
		case <-m.stopChan:
			m.cleaner.Stop()
			return
		}
	}
}

func (m *memoryCache) cleanExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for key, item := range m.items {
		if item.hasExpiry && now.After(item.expiration) {
			delete(m.items, key)
		}
	}
}

func (m *memoryCache) isExpired(item *cacheItem) bool {
	if !item.hasExpiry {
		return false
	}
	return time.Now().After(item.expiration)
}

func (m *memoryCache) IsExist(ctx context.Context, key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[m.key(key)]
	if !exists {
		return false
	}

	if m.isExpired(item) {
		// 延迟删除过期项
		go func() {
			m.mu.Lock()
			delete(m.items, m.key(key))
			m.mu.Unlock()
		}()
		return false
	}

	return true
}

func (m *memoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[m.key(key)]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}

	if m.isExpired(item) {
		// 延迟删除过期项
		go func() {
			m.mu.Lock()
			delete(m.items, m.key(key))
			m.mu.Unlock()
		}()
		return "", fmt.Errorf("key expired: %s", key)
	}

	return item.value, nil
}

func (m *memoryCache) Put(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := &cacheItem{
		value: value,
	}

	if expiration > 0 {
		item.expiration = time.Now().Add(expiration)
		item.hasExpiry = true
	}

	m.items[m.key(key)] = item
	return nil
}

func (m *memoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, m.key(key))
	return nil
}

func (m *memoryCache) BatchDelete(ctx context.Context, prefix string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	prefixKey := m.key(prefix)
	for key := range m.items {
		if strings.HasPrefix(key, prefixKey) {
			delete(m.items, key)
		}
	}

	return nil
}

func (m *memoryCache) GetWithTTL(ctx context.Context, key string) (string, time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[m.key(key)]
	if !exists {
		return "", 0, fmt.Errorf("key not found: %s", key)
	}

	if m.isExpired(item) {
		// 延迟删除过期项
		go func() {
			m.mu.Lock()
			delete(m.items, m.key(key))
			m.mu.Unlock()
		}()
		return "", 0, fmt.Errorf("key expired: %s", key)
	}

	var ttl time.Duration
	if item.hasExpiry {
		ttl = time.Until(item.expiration)
		if ttl < 0 {
			ttl = 0
		}
	} else {
		ttl = -1 // 表示永不过期
	}

	return item.value, ttl, nil
}

func (m *memoryCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[m.key(key)]
	if !exists {
		// 如果键不存在，创建一个新的
		item = &cacheItem{
			value: "0",
		}
		m.items[m.key(key)] = item
	} else if m.isExpired(item) {
		// 如果过期，重置为0
		item.value = "0"
		item.hasExpiry = false
	}

	// 解析当前值
	var currentValue int64
	if _, err := fmt.Sscanf(item.value, "%d", &currentValue); err != nil {
		return 0, fmt.Errorf("value is not a number: %s", item.value)
	}

	// 增加值
	newValue := currentValue + delta
	item.value = fmt.Sprintf("%d", newValue)

	return newValue, nil
}

func (m *memoryCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return m.Increment(ctx, key, -delta)
}

func (m *memoryCache) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[m.key(key)]
	if exists && !m.isExpired(item) {
		return false, nil // 键已存在且未过期
	}

	// 设置新值
	newItem := &cacheItem{
		value: value,
	}

	if expiration > 0 {
		newItem.expiration = time.Now().Add(expiration)
		newItem.hasExpiry = true
	}

	m.items[m.key(key)] = newItem
	return true, nil
}

func (m *memoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[m.key(key)]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	if m.isExpired(item) {
		delete(m.items, m.key(key))
		return fmt.Errorf("key expired: %s", key)
	}

	if expiration > 0 {
		item.expiration = time.Now().Add(expiration)
		item.hasExpiry = true
	} else {
		item.hasExpiry = false
	}

	return nil
}

func (m *memoryCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var keys []string
	patternKey := m.key(pattern)

	for key, item := range m.items {
		if m.isExpired(item) {
			continue
		}

		// 简单的模式匹配，支持 * 通配符
		if pattern == "*" || strings.Contains(key, strings.ReplaceAll(patternKey, "*", "")) {
			// 移除前缀返回原始键名
			originalKey := key
			if m.keyPrefix != "" {
				originalKey = strings.TrimPrefix(key, m.keyPrefix+":")
			}
			keys = append(keys, originalKey)
		}
	}

	return keys, nil
}

func (m *memoryCache) FlushAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]*cacheItem)
	return nil
}

func (m *memoryCache) Close() error {
	m.stopChan <- true
	return nil
}
