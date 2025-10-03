package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Save 保存用户
	Save(ctx context.Context, user *model.User) error

	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id uint) (*model.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// List 获取用户列表
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)

	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
}
