package service

import (
	"context"

	"github.com/dysodeng/app/internal/domain/user/model"
	"github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
)

// UserDomainService 用户领域服务
type UserDomainService interface {
	Info(ctx context.Context, userId uint64) (*model.User, error)
	ListUser(ctx context.Context, telephone, nickname, realName string, status uint8, page, pageSize int) ([]model.User, int64, error)
	CreateUser(ctx context.Context, userInfo *model.User) error
	UpdateUser(ctx context.Context, userInfo *model.User) error
	DeleteUser(ctx context.Context, userId uint64) error
}

type userDomainService struct {
	baseTraceSpanName string
	userRepo          repository.UserRepository
}

func NewUserDomainService(userRepo repository.UserRepository) UserDomainService {
	return &userDomainService{
		baseTraceSpanName: "domain.user.service.UserDomainService",
		userRepo:          userRepo,
	}
}

func (svc *userDomainService) Info(ctx context.Context, userId uint64) (*model.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".Info")
	defer span.End()
	user, err := svc.userRepo.Info(spanCtx, userId)
	if err != nil {
		logger.Error(spanCtx, "用户查询失败", logger.ErrorField(err))
		return nil, errors.New("用户查询失败")
	}
	if user == nil || user.ID <= 0 {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (svc *userDomainService) ListUser(ctx context.Context, telephone, nickname, realName string, status uint8, page, pageSize int) ([]model.User, int64, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".ListUser")
	defer span.End()

	query := repository.UserListQuery{
		Telephone: telephone,
		Nickname:  nickname,
		RealName:  realName,
		Status:    status,
		Page:      page,
		PageSize:  pageSize,
	}

	users, total, err := svc.userRepo.ListUser(spanCtx, query)
	if err != nil {
		logger.Error(spanCtx, "用户查询失败", logger.ErrorField(err))
		return nil, 0, errors.New("用户查询失败")
	}

	return users, total, nil
}

func (svc *userDomainService) CreateUser(ctx context.Context, userInfo *model.User) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CreateUser")
	defer span.End()

	if userInfo.Password == "" {
		return errors.New("密码不能为空")
	}
	password, err := helper.GeneratePassword(userInfo.Password)
	if err != nil {
		return err
	}
	userInfo.Password = password

	err = svc.userRepo.CreateUser(spanCtx, userInfo)
	if err != nil {
		logger.Error(spanCtx, "用户创建失败", logger.ErrorField(err))
		return errors.New("用户创建失败")
	}

	return nil
}

func (svc *userDomainService) UpdateUser(ctx context.Context, userInfo *model.User) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CreateUser")
	defer span.End()

	if userInfo.ID <= 0 {
		return errors.New("缺少用户ID")
	}
	if userInfo.Password != "" {
		password, err := helper.GeneratePassword(userInfo.Password)
		if err != nil {
			return err
		}
		userInfo.Password = password
	}

	err := svc.userRepo.UpdateUser(spanCtx, userInfo)
	if err != nil {
		logger.Error(spanCtx, "用户创建失败", logger.ErrorField(err))
		return errors.New("用户创建失败")
	}

	return nil
}

func (svc *userDomainService) DeleteUser(ctx context.Context, userId uint64) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".DeleteUser")
	defer span.End()

	err := svc.userRepo.DeleteUser(spanCtx, userId)
	if err != nil {
		logger.Error(spanCtx, "用户删除失败", logger.ErrorField(err))
		return errors.New("用户删除失败")
	}

	return nil
}
