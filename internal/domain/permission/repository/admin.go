package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/permission/model"
	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
)

// AdminRepository 管理员仓储
type AdminRepository interface {
	FindById(ctx context.Context, id uint64) (*model.Admin, error)
	FindByUsername(ctx context.Context, username sharedVO.Username) (*model.Admin, error)
	ExistsByUsername(ctx context.Context, username sharedVO.Username) (bool, error)
	Save(ctx context.Context, admin *model.Admin) error
	ChangePassword(ctx context.Context, id uint64, password sharedVO.Password) error
}
