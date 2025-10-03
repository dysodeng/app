package server

import "context"

type Server interface {
	IsEnabled() bool
	Name() string
	Addr() string
	Start() error
	Stop(ctx context.Context) error
}
