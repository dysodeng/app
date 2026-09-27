package config

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// WatchCallback 配置变更回调函数。回调不得修改已发布的 cfg。
type WatchCallback func(cfg *Config)

// Watcher 配置变更监听器
type Watcher struct {
	source    Source
	localPath string
	callback  WatchCallback

	mu      sync.Mutex
	cancel  context.CancelFunc
	done    chan struct{}
	stopped bool
}

// NewWatcher 创建配置变更监听器，localPath 是每次更新的本地配置基础。
func NewWatcher(source Source, localPath string, callback WatchCallback) *Watcher {
	return &Watcher{source: source, localPath: localPath, callback: callback}
}

// Start 启动配置监听，连接失败时在后台重试。
func (w *Watcher) Start(ctx context.Context) error {
	if w.source == nil || w.localPath == "" {
		return errors.New("配置监听缺少配置源或本地配置路径")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		return errors.New("配置监听器已停止")
	}
	if w.done != nil {
		return nil
	}
	watchCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.done = make(chan struct{})
	go w.run(watchCtx, w.done)
	return nil
}

// Stop 停止监听并关闭配置源；可以重复调用。
func (w *Watcher) Stop() error {
	w.mu.Lock()
	if w.stopped {
		w.mu.Unlock()
		return nil
	}
	w.stopped = true
	cancel, done := w.cancel, w.done
	w.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if done != nil {
		<-done
	}
	if closer, ok := w.source.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

func (w *Watcher) run(ctx context.Context, done chan struct{}) {
	defer close(done)
	for ctx.Err() == nil {
		ch, err := w.source.Watch(ctx)
		switch {
		case err != nil:
			slog.Warn("启动配置监听失败，稍后重试", "error", err)
		case ch == nil:
			slog.Warn("配置监听未返回更新通道，稍后重试")
		default:
			w.consume(ctx, ch)
			if ctx.Err() == nil {
				slog.Warn("配置监听已断开，稍后重试")
			}
		}
		if ctx.Err() != nil {
			return
		}
		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (w *Watcher) consume(ctx context.Context, ch <-chan []byte) {
	// 监听建立后再拉取当前值，补上启动或重连期间可能遗漏的变更。
	var retry <-chan time.Time
	var ticker *time.Ticker
	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
	}()
	startRetry := func() {
		if ticker == nil {
			ticker = time.NewTicker(time.Second)
			retry = ticker.C
		}
	}
	stopRetry := func() {
		if ticker != nil {
			ticker.Stop()
			ticker = nil
			retry = nil
		}
	}
	if !w.syncCurrent() {
		startRetry()
	}
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-ch:
			if !ok {
				return
			}
			// 事件可能早于上面的同步读取，重新读取当前值以免旧事件覆盖新配置。
			current, err := w.source.Load()
			if err != nil {
				slog.Warn("读取配置中心当前值失败，稍后重试", "error", err)
				startRetry()
				continue
			}
			if len(current) > 0 {
				data = current
			}
			stopRetry()
			w.apply(data)
		case <-retry:
			if w.syncCurrent() {
				stopRetry()
			}
		}
	}
}

func (w *Watcher) syncCurrent() bool {
	data, err := w.source.Load()
	if err != nil {
		slog.Warn("同步配置中心当前值失败，稍后重试", "error", err)
		return false
	}
	if len(data) > 0 {
		w.apply(data)
	}
	return true
}

func (w *Watcher) apply(data []byte) {
	_ = godotenv.Load()
	v, err := readLocal(w.localPath)
	if err != nil {
		slog.Error("读取本地配置失败，忽略本次更新", "error", err)
		return
	}
	if err := v.MergeConfig(bytes.NewReader(data)); err != nil {
		slog.Error("解析配置中心变更失败，忽略本次更新", "error", err)
		return
	}
	cfg, err := decodeConfig(v)
	if err != nil {
		slog.Error("合并配置变更失败，忽略本次更新", "error", err)
		return
	}
	if reflect.DeepEqual(Current(), cfg) {
		return
	}
	slog.Info("检测到配置变更，正在应用")
	Publish(cfg)
	if w.callback != nil {
		w.callback(cfg)
	}
}
