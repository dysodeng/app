package permission

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/dysodeng/app/internal/domain/permission/model"
	"github.com/dysodeng/app/internal/domain/permission/repository"
	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
	"github.com/dysodeng/app/internal/infrastructure/persistence/entity/permission"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	sharedModel "github.com/dysodeng/app/internal/infrastructure/shared/model"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

type adminRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewAdminRepository(txManager transactions.TransactionManager) repository.AdminRepository {
	return &adminRepository{
		txManager: txManager,
	}
}

func (repo *adminRepository) FindById(ctx context.Context, id uint64) (*model.Admin, error) {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".FindById")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx).Debug()

	var info permission.Admin
	if err := tx.Where("id = ?", id).First(&info).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return repo.adminFromModel(&info), nil
}

func (repo *adminRepository) FindByUsername(ctx context.Context, username sharedVO.Username) (*model.Admin, error) {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".FindByUsername")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx).Debug()

	var info permission.Admin
	if err := tx.Where("username = ?", username.Value()).First(&info).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return repo.adminFromModel(&info), nil
}

func (repo *adminRepository) ExistsByUsername(ctx context.Context, username sharedVO.Username) (bool, error) {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".FindByUsername")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx).Debug()

	var info permission.Admin
	if err := tx.Where("username = ?", username.Value()).Select("id").First(&info).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	}

	return info.ID > 0, nil
}

func (repo *adminRepository) Save(ctx context.Context, admin *model.Admin) error {
	return nil
}

func (repo *adminRepository) ChangePassword(ctx context.Context, id uint64, password sharedVO.Password) error {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".ChangePassword")
	defer span.End()
	tx := repo.txManager.GetTx(spanCtx).Model(&permission.Admin{}).Debug()
	return tx.Where("id=?", id).Update("safe_password", password.Value()).Error
}

func (repo *adminRepository) adminFromModel(admin *permission.Admin) *model.Admin {
	username, _ := sharedVO.NewUsername(admin.Username)
	password, _ := sharedVO.NewPasswordByHashText(admin.SafePassword)
	telephone, _ := sharedVO.NewTelephone(admin.Telephone)
	return &model.Admin{
		ID:           admin.ID,
		Username:     username,
		SafePassword: password,
		RealName:     admin.RealName,
		Telephone:    telephone,
		Remark:       admin.Remark,
		IsSuper:      sharedModel.BinaryStatusByUint(admin.IsSuper),
		Status:       sharedModel.BinaryStatusByUint(admin.Status),
	}
}
