package cache

import (
	"context"

	"golang.org/x/sync/singleflight"
)

type SFGroup[T any] struct {
	g singleflight.Group
}

func (s *SFGroup[T]) Do(ctx context.Context, key string, fn func() (T, error)) (T, error) {
	v, err, _ := s.g.Do(key, func() (interface{}, error) {
		return fn()
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return v.(T), nil
}

func (s *SFGroup[T]) Forget(key string) {
	s.g.Forget(key)
}
