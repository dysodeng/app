package user

import (
	"context"

	"github.com/dysodeng/app/internal/dal/model/user"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
)

// Dao 用户数据访问层
type Dao interface {
	Info(ctx context.Context, id uint64) (*user.User, error)
	ListUser(ctx context.Context, page, pageSize int, condition map[string]interface{}) ([]user.User, int64, error)
	CreateUser(ctx context.Context, userInfo user.User) (*user.User, error)
	UpdateUser(ctx context.Context, userInfo user.User) (*user.User, error)
	DeleteUser(ctx context.Context, userId uint64) error
}

type dao struct {
	baseTraceSpanName string
}

func NewUserDao() Dao {
	return &dao{
		baseTraceSpanName: "dal.dao.user.UserDao",
	}
}

func (d *dao) Info(ctx context.Context, id uint64) (*user.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, d.baseTraceSpanName+".Info")
	defer span.End()
	var userInfo user.User
	db.DB().WithContext(spanCtx).Debug().Where("id=?", id).First(&userInfo)
	return &userInfo, nil
}

func (d *dao) ListUser(ctx context.Context, page, pageSize int, condition map[string]interface{}) ([]user.User, int64, error) {
	query := db.DB().Debug().WithContext(ctx).Model(&user.User{})
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

func (d *dao) CreateUser(ctx context.Context, userInfo user.User) (*user.User, error) {
	err := db.DB().Debug().WithContext(ctx).Create(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (d *dao) UpdateUser(ctx context.Context, userInfo user.User) (*user.User, error) {
	err := db.DB().WithContext(ctx).Model(&user.User{}).Where("id=?", userInfo.ID).Updates(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (d *dao) DeleteUser(ctx context.Context, userId uint64) error {
	return db.DB().WithContext(ctx).Delete(&user.User{}, userId).Error
}
