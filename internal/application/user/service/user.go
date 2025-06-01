package service

import (
	"context"
	"time"

	"github.com/dysodeng/app/internal/application/user/dto/command"
	"github.com/dysodeng/app/internal/application/user/dto/query"
	"github.com/dysodeng/app/internal/application/user/dto/response"
	"github.com/dysodeng/app/internal/domain/user/model"
	"github.com/dysodeng/app/internal/domain/user/service"
	pkgModel "github.com/dysodeng/app/internal/pkg/model"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
)

// ApplicationService 用户应用服务
type ApplicationService interface {
	Info(ctx context.Context, id uint64) (*response.UserResponse, error)
	UserList(ctx context.Context, query query.UserListQuery) (*response.UserListResponse, error)
	CreateUser(ctx context.Context, cmd *command.UserCreateCommand) (*response.UserResponse, error)
}

type userApplication struct {
	baseTraceSpanName string
	userDomainService service.UserDomainService
}

func NewUserApplication(userDomainService service.UserDomainService) ApplicationService {
	return &userApplication{
		baseTraceSpanName: "application.user.ApplicationService",
		userDomainService: userDomainService,
	}
}

func (svc *userApplication) Info(ctx context.Context, id uint64) (*response.UserResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".Info")
	defer span.End()

	user, err := svc.userDomainService.Info(spanCtx, id)
	if err != nil {
		return nil, err
	}

	return response.FromDomainUser(user), nil
}

func (svc *userApplication) UserList(ctx context.Context, query query.UserListQuery) (*response.UserListResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".UserList")
	defer span.End()

	query.CheckOrDefault()

	list, total, err := svc.userDomainService.ListUser(spanCtx, query.Telephone, query.Nickname, query.RealName, query.Status, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	return &response.UserListResponse{
		List:  response.ListFromDomainUser(list),
		Total: total,
	}, nil
}

func (svc *userApplication) CreateUser(ctx context.Context, cmd *command.UserCreateCommand) (*response.UserResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CreateUser")
	defer span.End()

	var birthday time.Time
	if cmd.Birthday != "" {
		var err error
		birthday, err = time.ParseInLocation(time.DateOnly, cmd.Birthday, time.Local)
		if err != nil {
			return nil, errors.New("生日格式不正确")
		}
	}

	user := &model.User{
		Telephone: cmd.Telephone,
		Nickname:  cmd.Nickname,
		RealName:  cmd.RealName,
		Avatar:    cmd.Avatar,
		Gender:    cmd.Gender,
		Birthday:  birthday,
		Status:    pkgModel.BinaryStatusByUint(cmd.Status),
	}
	err := svc.userDomainService.CreateUser(spanCtx, user)
	if err != nil {
		return nil, err
	}

	return response.FromDomainUser(user), nil
}
