package user

import (
	"context"

	"github.com/dysodeng/app/internal/dal/model/user"
	"github.com/dysodeng/app/internal/pkg/db"
)

type Dao struct {
	ctx context.Context
}

func NewUserDao(ctx context.Context) *Dao {
	return &Dao{ctx: ctx}
}

func (d *Dao) Info(id uint64) (*user.User, error) {
	var userInfo user.User
	db.DB().WithContext(d.ctx).Where("id=?", id).First(&userInfo)
	return &userInfo, nil
}

func (d *Dao) ListUser(page, pageSize int, condition map[string]interface{}) ([]user.User, int64, error) {
	query := db.DB().Debug().WithContext(d.ctx).Model(&user.User{})
	if condition != nil {
		for where, val := range condition {
			query = query.Where(where, val)
		}
	}

	var userList []user.User
	var count int64
	query.Count(&count)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&userList)

	return userList, count, nil
}

func (d *Dao) CreateUser(userInfo user.User) (*user.User, error) {
	err := db.DB().Debug().WithContext(d.ctx).Create(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (d *Dao) UpdateUser(userInfo user.User) (*user.User, error) {
	err := db.DB().WithContext(d.ctx).Model(&user.User{}).Where("id=?", userInfo.ID).Updates(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (d *Dao) DeleteUser(userId uint64) error {
	return db.DB().WithContext(d.ctx).Delete(&user.User{}, userId).Error
}
