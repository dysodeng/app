package server

import "context"

// Server 服务器接口
type Server interface {
	// IsEnabled 是否启用
	IsEnabled() bool
	// Name 服务器名称
	Name() string
	// Addr 服务启动地址
	Addr() string
	// Start 启动服务
	Start() error
	// Stop 停止服务
	Stop(ctx context.Context) error
}
