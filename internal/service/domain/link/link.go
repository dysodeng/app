package link

import (
	"context"

	"github.com/dysodeng/app/internal/dal/model/link"
	"github.com/dysodeng/app/internal/pkg/vars"
	"github.com/pkg/errors"
)

// DomainService 跳转链接领域服务
type DomainService struct {
	ctx context.Context
}

func NewLinkDomainService(ctx context.Context) *DomainService {
	return &DomainService{ctx: ctx}
}

// Check 检查跳转链接
func (d *DomainService) Check(linkType link.Type, params link.Params) (*link.Link, error) {
	linkItem := link.Link{
		Type: linkType,
	}

	switch linkType {
	case link.Empty:
		linkItem.Params = link.Params{}

	case link.Navigation: // 导航
		if params.Point == nil {
			return nil, errors.New("缺少导航位置信息")
		}
		if vars.Float64Value(params.Point.Latitude) <= 0 || vars.Float64Value(params.Point.Longitude) <= 0 {
			return nil, errors.New("缺少导航位置信息")
		}
		linkItem.Params = link.Params{
			Point: params.Point,
		}

	case link.MiniProgram: // 小程序
		if vars.StringValue(params.AppID) == "" {
			return nil, errors.New("缺少小程序AppID")
		}
		if vars.StringValue(params.Path) == "" {
			return nil, errors.New("缺少小程序地址")
		}
		linkItem.Params = link.Params{
			AppID: params.AppID,
			Path:  params.Path,
		}

	default:
		return nil, errors.New("跳转链接类型错误")
	}

	return &linkItem, nil
}

// Build 构建可视化的跳转链接
func (d *DomainService) Build(linkItem link.Link) (*link.Link, error) {
	switch linkItem.Type {
	case link.Empty:
		break

	case link.Navigation:
		break

	case link.MiniProgram:
		break
	}
	return &linkItem, nil
}
