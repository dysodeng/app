package config

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentConfigPublication(t *testing.T) {
	previous := Current()
	t.Cleanup(func() { Publish(previous) })

	const updates = 1000
	var workers sync.WaitGroup
	workers.Add(2)
	start := make(chan struct{})
	go func() {
		defer workers.Done()
		<-start
		for i := 0; i < updates; i++ {
			version := fmt.Sprint(i)
			Publish(&Config{App: AppConfig{Name: version, Environment: version}})
		}
	}()
	go func() {
		defer workers.Done()
		<-start
		for i := 0; i < updates; i++ {
			if cfg := Current(); cfg != nil {
				if cfg.App.Name != cfg.App.Environment {
					t.Errorf("observed mixed config snapshot: %+v", cfg.App)
				}
			}
		}
	}()
	close(start)
	workers.Wait()
}
