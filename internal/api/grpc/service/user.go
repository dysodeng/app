package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dysodeng/app/internal/config"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/pkg/validator"
	userDo "github.com/dysodeng/app/internal/service/do/user"
	"github.com/dysodeng/app/internal/service/domain/user"
	"github.com/dysodeng/rpc/metadata"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// UserService 用户服务
type UserService struct {
	proto.UnimplementedUserServiceServer
	metadata.UnimplementedServiceRegister
}

func NewUserService() *UserService {
	return &UserService{}
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

	userDomainService := user.InitUserDomainService()
	userInfo, err := userDomainService.Info(spanCtx, req.Id)
	if err != nil {
		trace.Error(errors.Wrap(err, "获取用户信息失败"), span)
		logger.Error(spanCtx, "获取用户信息失败", logger.ErrorField(err))
		return nil, err
	}

	if userInfo.Id <= 0 {
		span.SetAttributes(attribute.String("query.user_id", fmt.Sprintf("%d", req.Id)))
		trace.Error(errors.New("用户不存在"), span)
		return nil, errors.New("用户不存在")
	} else {
		span.SetStatus(codes.Ok, "获取用户信息成功")
		span.SetAttributes(attribute.String("user_id", fmt.Sprintf("%d", userInfo.Id)))
		span.SetAttributes(attribute.String("nickname", userInfo.Nickname))
	}

	return &proto.UserResponse{
		Id:        userInfo.Id,
		Telephone: userInfo.Telephone,
		RealName:  userInfo.RealName,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Birthday:  userInfo.Birthday.Format(time.DateOnly),
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

	condition := map[string]interface{}{}
	if req.Username != "" {
		condition["username like %?%"] = req.Username
	}

	userDomainService := user.InitUserDomainService()
	list, total, err := userDomainService.ListUser(spanCtx, int(req.PageNum), int(req.PageSize), condition)
	if err != nil {
		logger.Error(spanCtx, "获取用户列表失败", logger.ErrorField(err))
		return nil, err
	}

	userList := make([]*proto.UserResponse, len(list))
	for i, item := range list {
		userList[i] = &proto.UserResponse{
			Id:        item.Id,
			Telephone: item.Telephone,
			RealName:  item.RealName,
			Nickname:  item.Nickname,
			Avatar:    item.Avatar,
			Birthday:  item.Birthday.Format(time.DateOnly),
			Gender:    uint32(item.Gender),
		}
	}

	return &proto.UserListResponse{
		List:  userList,
		Total: uint64(total),
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
	birthday, err := time.ParseInLocation(time.DateOnly, req.Birthday, time.Local)
	if err != nil {
		return nil, errors.New("生日格式不正确")
	}

	userDomainService := user.InitUserDomainService()
	userInfo, err := userDomainService.CreateUser(spanCtx, userDo.User{
		Telephone: req.Telephone,
		Password:  req.Password,
		RealName:  req.RealName,
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Birthday:  birthday,
		Gender:    uint8(req.Gender),
	})
	if err != nil {
		logger.Error(spanCtx, "创建用户失败", logger.ErrorField(err))
		return nil, err
	}
	return &proto.UserResponse{
		Id:        userInfo.Id,
		Telephone: userInfo.Telephone,
		RealName:  userInfo.RealName,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Birthday:  userInfo.Birthday.Format(time.DateOnly),
		Gender:    uint32(userInfo.Gender),
	}, nil
}
