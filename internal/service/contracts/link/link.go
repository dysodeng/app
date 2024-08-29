package link

import (
	"github.com/dysodeng/app/internal/model/link"
	"github.com/dysodeng/app/internal/service"
)

// ServiceInterface 跳转链接服务
type ServiceInterface interface {
	// Check 检查跳转链接
	Check(linkType link.Type, params link.Params) (*link.Link, service.Error)
	// Build 构建可视化的跳转链接
	Build(linkItem link.Link) (*link.Link, service.Error)
}
