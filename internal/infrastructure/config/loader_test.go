package config

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

type stubSource struct {
	data       []byte
	err        error
	updates    chan []byte
	closed     chan struct{}
	loadCalls  *atomic.Int32
	failLoads  int32
	watchCalls *atomic.Int32
}

func (s stubSource) Load() ([]byte, error) {
	if s.loadCalls != nil && s.loadCalls.Add(1) <= s.failLoads {
		return nil, errors.New("initial load failed")
	}
	return s.data, s.err
}

func (s stubSource) Type() string { return "test" }

func (s stubSource) Watch(context.Context) (<-chan []byte, error) {
	if s.watchCalls != nil && s.watchCalls.Add(1) == 1 {
		return nil, errors.New("initial watch failed")
	}
	return s.updates, nil
}

func (s stubSource) Close() error {
	if s.closed != nil {
		close(s.closed)
	}
	return nil
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	previous, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, previous)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}

func TestLoadPriorityEnvironmentRemoteLocal(t *testing.T) {
	previous := Current()
	t.Cleanup(func() { Publish(previous) })

	dir := t.TempDir()
	t.Chdir(dir)
	for _, key := range []string{"APP_ENV", "APP_DOMAIN", "SECURITY_JWT_SECRET", "CONFIG_CENTER_ETCD_KEY"} {
		unsetEnv(t, key)
	}
	local := []byte("config_center:\n  type: static\n  etcd:\n    key: local-key\napp:\n  name: local-name\n  environment: local-env\n  domain: local-domain\nsecurity:\n  jwt:\n    secret: local-secret\n")
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), local, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("APP_NAME=dotenv-name\nAPP_ENV=dotenv-env\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_NAME", "process-name")
	t.Setenv("CONFIG_CENTER_TYPE", "etcd")
	t.Setenv("CONFIG_CENTER_ETCD_ENDPOINTS", "127.0.0.1:2379,127.0.0.2:2379")
	t.Setenv("STORAGE_DRIVER", "env-storage")
	remote := []byte("app:\n  name: remote-name\n  environment: remote-env\nsecurity:\n  jwt:\n    secret: remote-secret\n")

	result, err := load("config.yaml", func(cc *ConfigCenterConfig) Source {
		if cc == nil || cc.Type != "etcd" || cc.Etcd == nil || cc.Etcd.Key != "local-key" ||
			len(cc.Etcd.Endpoints) != 2 || cc.Etcd.Endpoints[0] != "127.0.0.1:2379" {
			t.Errorf("unexpected bootstrap config: %+v", cc)
		}
		return stubSource{data: remote}
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Source != "test" || Current() != result.Config {
		t.Fatalf("remote configuration was not published: %+v", result)
	}
	if got := result.Config.App.Name; got != "process-name" {
		t.Errorf("process env should override remote and .env: got %q", got)
	}
	if got := result.Config.App.Environment; got != "dotenv-env" {
		t.Errorf(".env should override remote: got %q", got)
	}
	if got := result.Config.App.Domain; got != "local-domain" {
		t.Errorf("local field should survive remote omission: got %q", got)
	}
	if got := result.Config.Security.JWT.Secret; got != "remote-secret" {
		t.Errorf("remote field should override local: got %q", got)
	}
	if got := result.Config.Storage.Driver; got != "env-storage" {
		t.Errorf("environment should fill a section missing from files: got %q", got)
	}
}

func TestLoadFallsBackToLocalWhenRemoteFails(t *testing.T) {
	previous := Current()
	t.Cleanup(func() { Publish(previous) })
	unsetEnv(t, "APP_DOMAIN")

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: local-name\n  domain: local-domain\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_NAME", "env-name")

	result, err := load(path, func(*ConfigCenterConfig) Source {
		return stubSource{err: errors.New("remote unavailable")}
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Source != "local" || result.WatchSource == nil || result.Config.App.Name != "env-name" || result.Config.App.Domain != "local-domain" {
		t.Fatalf("unexpected fallback configuration: %+v", result)
	}
}

func TestWatcherStopsAndClosesSource(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: local-name\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	closed := make(chan struct{})
	watcher := NewWatcher(stubSource{updates: make(chan []byte), closed: closed}, path, nil)
	if err := watcher.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := watcher.Stop(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-closed:
	default:
		t.Fatal("watch source was not closed")
	}
	if err := watcher.Stop(); err != nil {
		t.Fatalf("second stop should be safe: %v", err)
	}
}

func TestWatcherRecoversAfterInitialLoadFailure(t *testing.T) {
	previous := Current()
	t.Cleanup(func() { Publish(previous) })
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: local-name\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	unsetEnv(t, "APP_NAME")
	loadCalls := &atomic.Int32{}
	source := stubSource{data: []byte("app:\n  name: remote-name\n"), updates: make(chan []byte), loadCalls: loadCalls, failLoads: 2}
	result, err := load(path, func(*ConfigCenterConfig) Source { return source })
	if err != nil {
		t.Fatal(err)
	}
	if result.Config.App.Name != "local-name" {
		t.Fatalf("initial load should use local config: %q", result.Config.App.Name)
	}
	applied := make(chan *Config, 1)
	watcher := NewWatcher(result.WatchSource, result.LocalPath, func(cfg *Config) { applied <- cfg })
	if err := watcher.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = watcher.Stop() })
	select {
	case cfg := <-applied:
		if cfg.App.Name != "remote-name" || Current() != cfg {
			t.Fatalf("recovered remote config was not published: %+v", cfg.App)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for recovered remote config")
	}
}

func TestWatcherRetriesFailedWatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: local-name\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	attempts := &atomic.Int32{}
	updates := make(chan []byte, 1)
	watcher := NewWatcher(stubSource{updates: updates, watchCalls: attempts}, path, nil)
	if err := watcher.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = watcher.Stop() })
	deadline := time.After(3 * time.Second)
	for attempts.Load() < 2 {
		select {
		case <-deadline:
			t.Fatal("watch was not retried")
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func TestWatcherDoesNotApplyQueuedStaleEventAfterSync(t *testing.T) {
	previous := Current()
	t.Cleanup(func() { Publish(previous) })
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: local-name\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	unsetEnv(t, "APP_NAME")
	updates := make(chan []byte, 1)
	updates <- []byte("app:\n  name: stale-name\n")
	loadCalls := &atomic.Int32{}
	watcher := NewWatcher(stubSource{data: []byte("app:\n  name: newest-name\n"), updates: updates, loadCalls: loadCalls}, path, nil)
	if err := watcher.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = watcher.Stop() })
	deadline := time.After(3 * time.Second)
	for loadCalls.Load() < 2 {
		select {
		case <-deadline:
			t.Fatal("queued event was not processed")
		case <-time.After(10 * time.Millisecond):
		}
	}
	if err := watcher.Stop(); err != nil {
		t.Fatal(err)
	}
	if cfg := Current(); cfg == nil || cfg.App.Name != "newest-name" {
		t.Fatalf("queued stale event rolled back configuration: %+v", cfg)
	}
}

func TestExplicitEmptyEnvironmentOverridesRemote(t *testing.T) {
	previous := Current()
	t.Cleanup(func() { Publish(previous) })

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: local-name\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_NAME", "")
	result, err := load(path, func(*ConfigCenterConfig) Source {
		return stubSource{data: []byte("app:\n  name: remote-name\n")}
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Config.App.Name != "" {
		t.Fatalf("explicit empty environment value should clear remote name, got %q", result.Config.App.Name)
	}
}

func TestWatcherReappliesPriorityForEveryRemoteUpdate(t *testing.T) {
	previous := Current()
	t.Cleanup(func() { Publish(previous) })
	unsetEnv(t, "APP_ENV")
	unsetEnv(t, "APP_DOMAIN")

	path := filepath.Join(t.TempDir(), "config.yaml")
	local := []byte("app:\n  name: local-name\n  environment: local-env\n  domain: local-domain\n")
	if err := os.WriteFile(path, local, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_NAME", "env-name")
	updates := make(chan []byte, 2)
	applied := make(chan *Config, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watcher := NewWatcher(stubSource{updates: updates}, path, func(cfg *Config) {
		applied <- cfg
	})
	if err := watcher.Start(ctx); err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		remote      string
		environment string
	}{
		{remote: "app:\n  name: remote-name\n  environment: remote-env\n", environment: "remote-env"},
		{remote: "app:\n  name: newer-name\n", environment: "local-env"},
	} {
		updates <- []byte(test.remote)
		select {
		case cfg := <-applied:
			if Current() != cfg || cfg.App.Name != "env-name" || cfg.App.Environment != test.environment || cfg.App.Domain != "local-domain" {
				t.Errorf("unexpected effective config after update: %+v", cfg.App)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("timed out waiting for config update")
		}
	}
}
