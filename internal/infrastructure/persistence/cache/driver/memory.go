package driver

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"
)

type memoryItem struct {
	val      []byte
	expireAt time.Time
}

// Memory 内存驱动
type Memory struct {
	mu       sync.RWMutex
	items    map[string]memoryItem
	gcTicker *time.Ticker
	gcStop   chan struct{}
}

func NewMemoryCache() *Memory {
	mb := &Memory{
		items:    make(map[string]memoryItem, 1024),
		gcTicker: time.NewTicker(time.Minute),
		gcStop:   make(chan struct{}),
	}
	go mb.gc()
	return mb
}

func (m *Memory) Exists(_ context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	it, ok := m.items[key]
	if !ok {
		return false, nil
	}
	if !it.expireAt.IsZero() && time.Now().After(it.expireAt) {
		return false, nil
	}
	return true, nil
}

func (m *Memory) Get(_ context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	it, ok := m.items[key]
	if !ok {
		return nil, nil
	}
	if !it.expireAt.IsZero() && time.Now().After(it.expireAt) {
		return nil, nil
	}
	return it.val, nil
}

func (m *Memory) Set(_ context.Context, key string, val []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	m.items[key] = memoryItem{val: val, expireAt: exp}
	return nil
}

func (m *Memory) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
	return nil
}

func (m *Memory) ScanDeleteByPrefix(_ context.Context, prefix string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.items {
		if strings.HasPrefix(k, prefix) {
			delete(m.items, k)
		}
	}
	return nil
}

func (m *Memory) Incr(_ context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	it, ok := m.items[key]
	var cur int64
	if ok && len(it.val) > 0 {
		if v, err := strconv.ParseInt(string(it.val), 10, 64); err == nil {
			cur = v
		}
	}
	cur++
	m.items[key] = memoryItem{val: []byte(strconv.FormatInt(cur, 10))}
	return cur, nil
}

func (m *Memory) gc() {
	for {
		select {
		case <-m.gcTicker.C:
			now := time.Now()
			m.mu.Lock()
			for k, it := range m.items {
				if !it.expireAt.IsZero() && now.After(it.expireAt) {
					delete(m.items, k)
				}
			}
			m.mu.Unlock()
		case <-m.gcStop:
			return
		}
	}
}
