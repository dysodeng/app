package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/model"
	domainService "github.com/dysodeng/app/internal/domain/service"
)

// UserAppService 用户应用服务
type UserAppService struct {
	userService *domainService.UserService
}

// NewUserAppService 创建用户应用服务
func NewUserAppService(userService *domainService.UserService) *UserAppService {
	return &UserAppService{
		userService: userService,
	}
}

// RegisterUser 注册用户
func (s *UserAppService) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	return s.userService.Register(ctx, username, email, password)
}

// GetUser 获取用户信息
func (s *UserAppService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userService.GetUserByID(ctx, id)
}

// GetUserList 获取用户列表
func (s *UserAppService) GetUserList(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return s.userService.GetUserList(ctx, page, pageSize)
}

// DeleteUser 删除用户
func (s *UserAppService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userService.DeleteUser(ctx, id)
}
