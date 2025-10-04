package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/model"
	"github.com/dysodeng/app/internal/domain/repository"
)

// UserService 用户领域服务
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户领域服务
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	// 检查用户名是否已存在
	existUser, _ := s.userRepo.FindByUsername(ctx, username)
	if existUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existUser, _ = s.userRepo.FindByEmail(ctx, email)
	if existUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 创建新用户
	user := model.NewUser(username, email, password)
	err := s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return s.userRepo.List(ctx, page, pageSize)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
