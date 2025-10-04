package repository

import (
	"context"

	"github.com/google/uuid"

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
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	if user.ID == uuid.Nil {
		if err := r.tx.GetTx(ctx).Create(userEntity).Error; err != nil {
			return err
		}
		user.ID = userEntity.ID
		user.CreatedAt = userEntity.CreatedAt.Time
		user.UpdatedAt = userEntity.UpdatedAt.Time
	} else {
		if err := r.tx.GetTx(ctx).Save(userEntity).Error; err != nil {
			return err
		}
		user.UpdatedAt = userEntity.UpdatedAt.Time
	}

	return nil
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var userEntity entity.User
	result := r.tx.GetTx(ctx).Debug().First(&userEntity, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.User{
		ID:        userEntity.ID,
		Username:  userEntity.Username,
		Email:     userEntity.Email,
		Password:  userEntity.Password,
		CreatedAt: userEntity.CreatedAt.Time,
		UpdatedAt: userEntity.UpdatedAt.Time,
	}, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var userEntity entity.User
	result := r.tx.GetTx(ctx).Debug().Where("username = ?", username).First(&userEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.User{
		ID:        userEntity.ID,
		Username:  userEntity.Username,
		Email:     userEntity.Email,
		Password:  userEntity.Password,
		CreatedAt: userEntity.CreatedAt.Time,
		UpdatedAt: userEntity.UpdatedAt.Time,
	}, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var userEntity entity.User
	result := r.tx.GetTx(ctx).Debug().Where("email = ?", email).First(&userEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.User{
		ID:        userEntity.ID,
		Username:  userEntity.Username,
		Email:     userEntity.Email,
		Password:  userEntity.Password,
		CreatedAt: userEntity.CreatedAt.Time,
		UpdatedAt: userEntity.UpdatedAt.Time,
	}, nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var userEntities []entity.User
	var total int64

	offset := (page - 1) * pageSize
	result := r.tx.GetTx(ctx).Debug().Model(&entity.User{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = r.tx.GetTx(ctx).Debug().Offset(offset).Limit(pageSize).Find(&userEntities)
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
			CreatedAt: userEntity.CreatedAt.Time,
			UpdatedAt: userEntity.UpdatedAt.Time,
		}
	}

	return users, total, nil
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.tx.GetTx(ctx).Debug().Delete(&entity.User{}, id)
	return result.Error
}
