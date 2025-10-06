package model

import (
	"time"

	"github.com/google/uuid"

	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
	"github.com/dysodeng/app/internal/domain/user/valueobject"
	sharedModel "github.com/dysodeng/app/internal/infrastructure/shared/model"
)

// User 用户领域模型
type User struct {
	ID                  uuid.UUID
	Telephone           sharedVO.Telephone
	WxUnionID           valueobject.WxUnionID
	WxMiniProgramOpenID valueobject.WxMiniProgramOpenID
	WxOfficialOpenID    valueobject.WxOfficialOpenID
	Nickname            string
	Avatar              valueobject.Avatar
	Status              sharedModel.BinaryStatus
	CreatedAt           time.Time
}

func NewUser(
	telephone sharedVO.Telephone,
	wxUnionID valueobject.WxUnionID,
	wxMiniProgramOpenID valueobject.WxMiniProgramOpenID,
	avatar valueobject.Avatar,
	nickname string,
) (*User, error) {
	id, _ := uuid.NewV7()
	u := &User{
		ID:                  id,
		Telephone:           telephone,
		WxUnionID:           wxUnionID,
		WxMiniProgramOpenID: wxMiniProgramOpenID,
		Nickname:            nickname,
		Avatar:              avatar,
		Status:              sharedModel.BinaryStatusTrue,
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) Validate() error {
	if err := u.Telephone.Validate(); err != nil {
		return err
	}
	if err := u.WxUnionID.Validate(); err != nil {
		return err
	}
	if err := u.WxMiniProgramOpenID.Validate(); err != nil {
		return err
	}
	if err := u.Avatar.Validate(); err != nil {
		return err
	}
	return nil
}
