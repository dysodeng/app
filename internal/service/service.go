package service

import (
	"github.com/defval/di"
	"github.com/dysodeng/app/internal/pkg/api"
)

// Error 服务错误码
type Error struct {
	Code  api.Code
	Error error
}

// Container 服务容器
var Container *di.Container

// Resolve 获取服务实例
func Resolve(ptr di.Pointer, options ...di.ResolveOption) error {
	return Container.Resolve(ptr, options...)
}
