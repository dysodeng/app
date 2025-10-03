package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/model"
	"github.com/dysodeng/app/internal/domain/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/entity"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	tx transactions.TransactionManager
}

// NewUserRepository 创建用户仓储
func NewUserRepository(tx transactions.TransactionManager) repository.UserRepository {
	return &UserRepositoryImpl{
		tx: tx,
	}
}

// Save 保存用户
func (r *UserRepositoryImpl) Save(ctx context.Context, user *model.User) error {
	userEntity := &entity.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	result := r.tx.GetTx(ctx).Save(userEntity)
	if result.Error != nil {
		return result.Error
	}

	user.ID = userEntity.ID
	return nil
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var userEntity entity.User
	result := r.tx.GetTx(ctx).First(&userEntity, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.User{
		ID:        userEntity.ID,
		Username:  userEntity.Username,
		Email:     userEntity.Email,
		Password:  userEntity.Password,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var userEntity entity.User
	result := r.tx.GetTx(ctx).Where("username = ?", username).First(&userEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.User{
		ID:        userEntity.ID,
		Username:  userEntity.Username,
		Email:     userEntity.Email,
		Password:  userEntity.Password,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var userEntity entity.User
	result := r.tx.GetTx(ctx).Where("email = ?", email).First(&userEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.User{
		ID:        userEntity.ID,
		Username:  userEntity.Username,
		Email:     userEntity.Email,
		Password:  userEntity.Password,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}, nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var userEntities []entity.User
	var total int64

	offset := (page - 1) * pageSize
	result := r.tx.GetTx(ctx).Model(&entity.User{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = r.tx.GetTx(ctx).Offset(offset).Limit(pageSize).Find(&userEntities)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	users := make([]*model.User, len(userEntities))
	for i, userEntity := range userEntities {
		users[i] = &model.User{
			ID:        userEntity.ID,
			Username:  userEntity.Username,
			Email:     userEntity.Email,
			Password:  userEntity.Password,
			CreatedAt: userEntity.CreatedAt,
			UpdatedAt: userEntity.UpdatedAt,
		}
	}

	return users, total, nil
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.tx.GetTx(ctx).Delete(&entity.User{}, id)
	return result.Error
}
