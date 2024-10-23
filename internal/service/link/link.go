package link

import (
	"github.com/dysodeng/app/internal/dal/model/link"
	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/vars"
	"github.com/dysodeng/app/internal/service"
	linkIface "github.com/dysodeng/app/internal/service/contracts/link"

	"github.com/pkg/errors"
)

// Service 跳转链接服务
type Service struct{}

var _ linkIface.ServiceInterface = (*Service)(nil)

func NewLinkService() *Service {
	return &Service{}
}

// Check 检查跳转链接
func (linkService Service) Check(linkType link.Type, params link.Params) (*link.Link, service.Error) {
	linkItem := link.Link{
		Type: linkType,
	}

	switch linkType {
	case link.Empty:
		linkItem.Params = link.Params{}
		break

	case link.Navigation: // 外部导航
		if params.Point == nil {
			return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少导航地址信息")}
		}
		if vars.Float64Value(params.Point.Latitude) <= 0 || vars.Float64Value(params.Point.Longitude) <= 0 {
			return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少导航地址信息")}
		}
		linkItem.Params = link.Params{
			Point: params.Point,
		}
		break

	case link.MiniProgram: // 跳转小程序
		if vars.StringValue(params.AppID) == "" {
			return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少小程序AppID")}
		}
		if vars.StringValue(params.Path) == "" {
			return nil, service.Error{Code: api.CodeFail, Error: errors.New("缺少小程序地址")}
		}
		linkItem.Params = link.Params{
			AppID: params.AppID,
			Path:  params.Path,
		}

	default:
		return nil, service.Error{Code: api.CodeFail, Error: errors.New("链接类型错误")}
	}

	return &linkItem, service.Error{Code: api.CodeOk}
}

// Build 构建可视化的跳转链接
func (linkService Service) Build(linkItem link.Link) (*link.Link, service.Error) {

	switch linkItem.Type {
	case link.Empty:
		break

	case link.Navigation: // 外部导航
		break

	case link.MiniProgram:
		break
	}

	return &linkItem, service.Error{Code: api.CodeOk}
}
