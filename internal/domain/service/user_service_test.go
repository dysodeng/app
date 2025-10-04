package service_test

import (
	"context"
	"testing"

	"github.com/dysodeng/app/internal/domain/model"
	domainService "github.com/dysodeng/app/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*model.User), 0, args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	// 创建模拟仓储
	mockRepo := new(MockUserRepository)

	// 设置模拟行为
	mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(nil, nil)
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	// 创建用户服务
	userService := domainService.NewUserService(mockRepo)

	// 执行注册方法
	user, err := userService.Register(context.Background(), "testuser", "test@example.com", "password")

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID(t *testing.T) {
	// 创建模拟仓储
	mockRepo := new(MockUserRepository)

	// 创建模拟用户
	id, _ := uuid.NewV7()
	mockUser := &model.User{
		ID:       id,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// 设置模拟行为
	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(mockUser, nil)

	// 创建用户服务
	userService := domainService.NewUserService(mockRepo)

	// 执行获取用户方法
	user, err := userService.GetUserByID(context.Background(), id)

	// 断言结果
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}
