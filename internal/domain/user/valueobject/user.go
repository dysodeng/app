package valueobject

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
)

const (
	DefaultUserAvatar = "static/avatar.png" // 默认用户头像
)

// WxUnionID 微信开放平台union_id值对象
type WxUnionID struct {
	value string
}

func NewWxUnionID(unionId string) (WxUnionID, error) {
	return WxUnionID{
		value: unionId,
	}, nil
}

func (u WxUnionID) String() string {
	return u.Value()
}

func (u WxUnionID) Value() string {
	return u.value
}

func (u WxUnionID) Validate() error {
	return nil
}

// WxMiniProgramOpenID 微信小程序用户OpenID值对象
type WxMiniProgramOpenID struct {
	value string
}

func NewWxMiniProgramOpenID(openid string) (WxMiniProgramOpenID, error) {
	openid = strings.TrimSpace(openid)
	vo := WxMiniProgramOpenID{value: openid}
	if err := vo.Validate(); err != nil {
		return WxMiniProgramOpenID{}, err
	}
	return vo, nil
}

func (o WxMiniProgramOpenID) String() string {
	return o.Value()
}

func (o WxMiniProgramOpenID) Value() string {
	return o.value
}

func (o WxMiniProgramOpenID) Validate() error {
	if o.value == "" {
		return errors.New("OpenID为空")
	}
	return nil
}

// WxOfficialOpenID 微信公众号用户OpenID值对象
type WxOfficialOpenID struct {
	value string
}

func NewWxOfficialOpenID(openid string) (WxOfficialOpenID, error) {
	openid = strings.TrimSpace(openid)
	vo := WxOfficialOpenID{value: openid}
	if err := vo.Validate(); err != nil {
		return WxOfficialOpenID{}, err
	}
	return vo, nil
}

func (o WxOfficialOpenID) String() string {
	return o.Value()
}

func (o WxOfficialOpenID) Value() string {
	return o.value
}

func (o WxOfficialOpenID) Validate() error {
	if o.value == "" {
		return errors.New("OpenID为空")
	}
	return nil
}

// Avatar 头像值对象
type Avatar struct {
	value string
}

func NewAvatar(avatar string) (Avatar, error) {
	ctx := context.Background()
	avatar = strings.TrimSpace(avatar)
	if avatar == "" {
		avatar = DefaultUserAvatar
	}
	return Avatar{
		value: storage.Instance().FullUrl(ctx, storage.Instance().RelativePath(ctx, avatar)),
	}, nil
}

func (a Avatar) String() string {
	return a.Value()
}

func (a Avatar) Value() string {
	return a.value
}

func (a Avatar) FullURL() string {
	return a.Value()
}

func (a Avatar) RelativePath() string {
	ctx := context.Background()
	return storage.Instance().RelativePath(ctx, a.value)
}

func (a Avatar) Validate() error {
	if a.value == "" {
		return errors.New("头像地址为空")
	}
	return nil
}
