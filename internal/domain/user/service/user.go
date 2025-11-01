package service

import (
	"context"

	"github.com/google/uuid"

	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
	"github.com/dysodeng/app/internal/domain/user/errors"
	"github.com/dysodeng/app/internal/domain/user/model"
	"github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/domain/user/valueobject"
)

type UserDomainService interface {
	UserInfo(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByTelephone(ctx context.Context, telephone string) (*model.User, error)
	FindByWxUnionId(ctx context.Context, wxUnionId string) (*model.User, error)
	FindByOpenId(ctx context.Context, platform, openid string) (*model.User, error)
	Create(ctx context.Context, telephone, unionId, wxMiniProgramOpenId, nickname, avatar string) (*model.User, error)
}

type userDomainService struct {
	userRepository repository.UserRepository
}

func NewUserDomainService(userRepository repository.UserRepository) UserDomainService {
	return &userDomainService{
		userRepository: userRepository,
	}
}

func (svc *userDomainService) UserInfo(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := svc.userRepository.FindById(ctx, id)
	if err != nil {
		return nil, errors.ErrUserQueryFailed
	}
	if user == nil || user.ID == uuid.Nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (svc *userDomainService) FindByTelephone(ctx context.Context, telephone string) (*model.User, error) {
	user, err := svc.userRepository.FindByTelephone(ctx, telephone)
	if err != nil {
		return nil, errors.ErrUserQueryFailed
	}
	if user == nil || user.ID == uuid.Nil {
		return &model.User{}, nil
	}
	return user, nil
}

func (svc *userDomainService) FindByWxUnionId(ctx context.Context, wxUnionId string) (*model.User, error) {
	user, err := svc.userRepository.FindByUnionId(ctx, wxUnionId)
	if err != nil {
		return nil, errors.ErrUserQueryFailed
	}
	if user == nil || user.ID == uuid.Nil {
		return &model.User{}, nil
	}
	return user, nil
}

func (svc *userDomainService) FindByOpenId(ctx context.Context, platform, openid string) (*model.User, error) {
	user, err := svc.userRepository.FindByOpenId(ctx, platform, openid)
	if err != nil {
		return nil, errors.ErrUserQueryFailed
	}
	if user == nil || user.ID == uuid.Nil {
		return &model.User{}, nil
	}
	return user, nil
}

// Create 创建用户
func (svc *userDomainService) Create(ctx context.Context, telephone, unionId, wxMiniProgramOpenId, nickname, avatar string) (*model.User, error) {
	u, err := svc.FindByTelephone(ctx, telephone)
	if err != nil {
		return nil, err
	}
	if u != nil && u.ID != uuid.Nil {
		return u, nil
	}

	telephoneVo, err := sharedVO.NewTelephone(telephone)
	if err != nil {
		return nil, err
	}
	unionIdVo, _ := valueobject.NewWxUnionID(unionId)
	wxMiniProgramOpenIdVo, err := valueobject.NewWxMiniProgramOpenID(wxMiniProgramOpenId)
	if err != nil {
		return nil, err
	}
	avatarVo, err := valueobject.NewAvatar(avatar)
	if err != nil {
		return nil, err
	}

	user, err := model.NewUser(telephoneVo, unionIdVo, wxMiniProgramOpenIdVo, avatarVo, nickname)
	if err != nil {
		return nil, err
	}
	if err = user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}
