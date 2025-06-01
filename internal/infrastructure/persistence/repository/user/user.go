package user

import (
	"context"
	"strings"

	"github.com/dysodeng/app/internal/domain/user/model"
	"github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/model/user"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
)

type userRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewUserRepository(txManager transactions.TransactionManager) repository.UserRepository {
	return &userRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.user.UserRepository",
		txManager:         txManager,
	}
}

func (repo *userRepository) Info(ctx context.Context, id uint64) (*model.User, error) {
	var userInfo user.User
	repo.txManager.GetTx(ctx).Debug().Where("id=?", id).First(&userInfo)
	return model.UserFromModel(&userInfo), nil
}

func (repo *userRepository) ListUser(ctx context.Context, query repository.UserListQuery) ([]model.User, int64, error) {
	tx := repo.txManager.GetTx(ctx)
	tx.Debug().Model(&user.User{})

	if query.Telephone != "" {
		if len(query.Telephone) == 11 {
			tx = tx.Where("telephone = ?", query.Telephone)
		} else {
			tx = tx.Where("name LIKE ?", query.Telephone+"%")
		}
	}
	if query.RealName != "" {
		tx = tx.Where("real_name LIKE ?", "%"+query.RealName+"%")
	}
	if query.Nickname != "" {
		tx = tx.Where("nickname LIKE ?", "%"+query.Nickname+"%")
	}
	if query.Status > 0 {
		tx = tx.Where("status = ?", query.Status-1)
	}

	// 获取总数
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if query.OrderBy != "" {
		orderType := "asc"
		if strings.ToLower(query.OrderType) == "desc" {
			orderType = "desc"
		}
		tx = tx.Order(query.OrderBy + " " + orderType)
	}

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		tx = tx.Offset(offset).Limit(query.PageSize)
	}

	var users []user.User
	if err := tx.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return model.UserListFromModel(users), total, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, userInfo *model.User) error {
	dataModel := userInfo.ToModel()
	dataModel.ID = 0
	err := repo.txManager.GetTx(ctx).Debug().Create(&dataModel).Error
	if err != nil {
		return err
	}
	userInfo.ID = dataModel.ID
	return nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, userInfo *model.User) error {
	dataModel := userInfo.ToModel()
	err := repo.txManager.GetTx(ctx).Debug().Model(&user.User{}).Where("id=?", userInfo.ID).Updates(&dataModel).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) DeleteUser(ctx context.Context, userId uint64) error {
	return repo.txManager.GetTx(ctx).Debug().Where("id=?", userId).Delete(&user.User{}).Error
}
