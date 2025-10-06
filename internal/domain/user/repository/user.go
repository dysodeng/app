package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/user/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByTelephone(ctx context.Context, telephone string) (*model.User, error)
	FindByUnionId(ctx context.Context, unionId string) (*model.User, error)
	FindByOpenId(ctx context.Context, platform, openId string) (*model.User, error)
	Save(ctx context.Context, userInfo *model.User) error
}
