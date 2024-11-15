package service

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/trace"
	"github.com/dysodeng/app/internal/service/rpc"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/validator"
	userDo "github.com/dysodeng/app/internal/service/do/user"
	"github.com/dysodeng/app/internal/service/domain/user"
	"github.com/dysodeng/rpc/metadata"
	"github.com/pkg/errors"
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
		ServiceName: "user.UserService",
		Version:     metadata.DefaultVersion,
	}
}

func (m *UserService) Info(ctx context.Context, req *proto.UserInfoRequest) (*proto.UserResponse, error) {
	traceCtx := trace.New().NewSpan(rpc.FromCtx(ctx), "grpc.user.Info")
	if req.Id <= 0 {
		return nil, errors.New("缺少用户ID")
	}

	userDomainService := user.NewUserDomainService(traceCtx)
	userInfo, err := userDomainService.Info(req.Id)
	if err != nil {
		logger.Error(traceCtx, "获取用户信息失败", logger.ErrorField(err))
		return nil, err
	}

	return &proto.UserResponse{
		Id:        userInfo.Id,
		Telephone: userInfo.Telephone,
		RealName:  userInfo.RealName,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Birthday:  userInfo.Birthday,
		Gender:    uint32(userInfo.Gender),
	}, nil
}

func (m *UserService) ListUser(ctx context.Context, req *proto.UserListRequest) (*proto.UserListResponse, error) {
	traceCtx := trace.New().NewSpan(rpc.FromCtx(ctx), "grpc.user.ListUser")
	logger.Debug(traceCtx, "获取用户列表接口", logger.Field{Key: "params", Value: req})
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

	userDomainService := user.NewUserDomainService(traceCtx)
	list, total, err := userDomainService.ListUser(int(req.PageNum), int(req.PageSize), condition)
	if err != nil {
		logger.Error(traceCtx, "获取用户列表失败", logger.ErrorField(err))
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
			Birthday:  item.Birthday,
			Gender:    uint32(item.Gender),
		}
	}

	return &proto.UserListResponse{
		List:  userList,
		Total: uint64(total),
	}, nil
}

func (m *UserService) CreateUser(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {
	traceCtx := trace.New().NewSpan(rpc.FromCtx(ctx), "grpc.user.CreateUser")
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
	userDomainService := user.NewUserDomainService(traceCtx)
	userInfo, err := userDomainService.CreateUser(userDo.User{
		Telephone: req.Telephone,
		Password:  req.Password,
		RealName:  req.RealName,
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Birthday:  req.Birthday,
		Gender:    uint8(req.Gender),
	})
	if err != nil {
		logger.Error(traceCtx, "创建用户失败", logger.ErrorField(err))
		return nil, err
	}
	return &proto.UserResponse{
		Id:        userInfo.Id,
		Telephone: userInfo.Telephone,
		RealName:  userInfo.RealName,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Birthday:  userInfo.Birthday,
		Gender:    uint32(userInfo.Gender),
	}, nil
}
