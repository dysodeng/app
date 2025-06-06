package service

import (
	"context"
	"fmt"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/application/user/dto/command"
	"github.com/dysodeng/app/internal/application/user/dto/query"
	"github.com/dysodeng/app/internal/application/user/service"
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/form"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/pkg/validator"
	"github.com/dysodeng/rpc/metadata"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// UserService 用户服务
type UserService struct {
	userApplicationService service.UserApplicationService
	proto.UnimplementedUserServiceServer
	metadata.UnimplementedServiceRegister
}

func NewUserService(userApplicationService service.UserApplicationService) *UserService {
	return &UserService{
		userApplicationService: userApplicationService,
	}
}

func (m *UserService) RegisterMetadata() metadata.ServiceRegisterMetadata {
	return metadata.ServiceRegisterMetadata{
		AppName:     config.App.Name,
		ServiceName: "user.UserService",
		Version:     metadata.DefaultVersion,
		Env:         config.App.Env.String(),
	}
}

func (m *UserService) Info(ctx context.Context, req *proto.UserInfoRequest) (*proto.UserResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, "grpc.user.Info")
	defer span.End()

	if req.Id <= 0 {
		return nil, errors.New("缺少用户ID")
	}

	userInfo, err := m.userApplicationService.Info(spanCtx, req.Id)
	if err != nil {
		trace.Error(errors.Wrap(err, "获取用户信息失败"), span)
		logger.Error(spanCtx, "获取用户信息失败", logger.ErrorField(err))
		return nil, err
	}

	if userInfo.ID <= 0 {
		span.SetAttributes(attribute.String("query.user_id", fmt.Sprintf("%d", req.Id)))
		trace.Error(errors.New("用户不存在"), span)
		return nil, errors.New("用户不存在")
	} else {
		span.SetStatus(codes.Ok, "获取用户信息成功")
		span.SetAttributes(attribute.String("user_id", fmt.Sprintf("%d", userInfo.ID)))
		span.SetAttributes(attribute.String("nickname", userInfo.Nickname))
	}

	return &proto.UserResponse{
		Id:        userInfo.ID,
		Telephone: userInfo.Telephone,
		RealName:  userInfo.RealName,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Birthday:  userInfo.Birthday,
		Gender:    uint32(userInfo.Gender),
	}, nil
}

func (m *UserService) ListUser(ctx context.Context, req *proto.UserListRequest) (*proto.UserListResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, "grpc.user.ListUser")
	defer span.End()
	logger.Debug(spanCtx, "获取用户列表接口", logger.Field{Key: "params", Value: req})
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	res, err := m.userApplicationService.UserList(spanCtx, query.UserListQuery{
		Pagination: form.Pagination{
			Page:     int(req.PageNum),
			PageSize: int(req.PageSize),
		},
		Telephone: req.Username,
	})
	if err != nil {
		logger.Error(spanCtx, "获取用户列表失败", logger.ErrorField(err))
		return nil, err
	}

	userList := make([]*proto.UserResponse, len(res.List))
	for i, item := range res.List {
		userList[i] = &proto.UserResponse{
			Id:        item.ID,
			Telephone: item.Telephone,
			RealName:  item.RealName,
			Nickname:  item.Nickname,
			Avatar:    item.Avatar,
			Birthday:  item.Birthday,
			Gender:    uint32(item.Gender),
		}
	}

	return &proto.UserListResponse{
		List:  userList,
		Total: uint64(res.Total),
	}, nil
}

func (m *UserService) CreateUser(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, "grpc.user.CreateUser")
	defer span.End()

	if req.Telephone == "" {
		return nil, errors.New("缺少手机号码")
	}
	if req.Password == "" {
		return nil, errors.New("登录缺少密码")
	}
	if req.RealName == "" {
		return nil, errors.New("缺少真实姓名")
	}
	if req.Nickname == "" {
		return nil, errors.New("缺少昵称")
	}
	if req.Avatar == "" {
		return nil, errors.New("缺少头像")
	}
	if !validator.IsMobile(req.Telephone) {
		return nil, errors.New("手机号码格式不正确")
	}
	if !validator.IsSafePassword(req.Password, 8) {
		return nil, errors.New("密码格式不正确")
	}
	if req.Birthday == "" {
		return nil, errors.New("缺少生日")
	}

	userInfo := &command.UserCreateCommand{
		Telephone: req.Telephone,
		Password:  req.Password,
		RealName:  req.RealName,
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Birthday:  req.Birthday,
		Gender:    uint8(req.Gender),
	}
	res, err := m.userApplicationService.CreateUser(spanCtx, userInfo)
	if err != nil {
		logger.Error(spanCtx, "创建用户失败", logger.ErrorField(err))
		return nil, err
	}

	return &proto.UserResponse{
		Id:        res.ID,
		Telephone: res.Telephone,
		RealName:  res.RealName,
		Nickname:  res.Nickname,
		Avatar:    res.Avatar,
		Birthday:  res.Birthday,
		Gender:    uint32(userInfo.Gender),
	}, nil
}
