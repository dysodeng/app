package model

import (
	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
	sharedModel "github.com/dysodeng/app/internal/infrastructure/shared/model"
)

type Admin struct {
	ID           uint64
	Username     sharedVO.Username
	SafePassword sharedVO.Password
	RealName     string
	Telephone    sharedVO.Telephone
	Remark       string
	IsSuper      sharedModel.BinaryStatus
	Status       sharedModel.BinaryStatus
}

func NewAdmin(
	username sharedVO.Username,
	password sharedVO.Password,
	realName string,
	telephone sharedVO.Telephone,
	remark string,
	status sharedModel.BinaryStatus,
) (*Admin, error) {
	return &Admin{
		Username:     username,
		SafePassword: password,
		RealName:     realName,
		Telephone:    telephone,
		Remark:       remark,
		IsSuper:      sharedModel.BinaryStatusFalse,
		Status:       status,
	}, nil
}
