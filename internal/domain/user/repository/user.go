package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/user/model"
)

// UserListQuery 用户列表查询条件
type UserListQuery struct {
	Telephone string
	RealName  string
	Nickname  string
	Status    uint8
	OrderBy   string // 排序字段
	OrderType string // 排序方式：asc/desc
	Page      int    // 页码
	PageSize  int    // 每页数量
}

type UserRepository interface {
	Info(ctx context.Context, id uint64) (*model.User, error)
	ListUser(ctx context.Context, query UserListQuery) ([]model.User, int64, error)
	CreateUser(ctx context.Context, userInfo *model.User) error
	UpdateUser(ctx context.Context, userInfo *model.User) error
	DeleteUser(ctx context.Context, userId uint64) error
}
