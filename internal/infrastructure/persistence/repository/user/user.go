package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
	"github.com/dysodeng/app/internal/domain/user/model"
	"github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/domain/user/valueobject"
	"github.com/dysodeng/app/internal/infrastructure/persistence/entity/user"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	sharedModel "github.com/dysodeng/app/internal/infrastructure/shared/model"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
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

func (repo *userRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".FindById")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx).Debug()

	var info user.User
	if err := tx.Where("id = ?", id).First(&info).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return repo.userFromModel(&info), nil
}

func (repo *userRepository) FindByTelephone(ctx context.Context, telephone string) (*model.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".FindByTelephone")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx).Debug()

	var info user.User
	if err := tx.Where("telephone = ?", telephone).First(&info).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return repo.userFromModel(&info), nil
}

func (repo *userRepository) FindByUnionId(ctx context.Context, unionId string) (*model.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".FindByUnionId")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx).Debug()

	var info user.User
	if err := tx.Where("union_id = ?", unionId).First(&info).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return repo.userFromModel(&info), nil
}

func (repo *userRepository) FindByOpenId(ctx context.Context, platform, openId string) (*model.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".FindByUnionId")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx).Debug()

	if platform == "WxMinioProgram" {
		tx = tx.Where("wx_mini_program_openid = ?", openId)
	} else {
		tx = tx.Where("wx_official_openid = ?", openId)
	}

	var info user.User
	if err := tx.First(&info).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return repo.userFromModel(&info), nil
}

func (repo *userRepository) Save(ctx context.Context, userInfo *model.User) error {
	spanCtx, span := trace.Tracer().Start(ctx, repo.baseTraceSpanName+".Save")
	defer span.End()

	tx := repo.txManager.GetTx(spanCtx)

	userModel := repo.toModel(userInfo)

	if userInfo.ID != uuid.Nil {
		var existsUser user.User
		tx.Where("id = ?", userModel.ID).First(&existsUser)
		if existsUser.ID == uuid.Nil {
			if err := tx.Create(userModel).Error; err != nil {
				return err
			}
			userInfo.ID = userModel.ID
			userInfo.CreatedAt = userModel.CreatedAt.Time
		} else {
			if err := tx.Where("id=?", userInfo.ID).
				Updates(userModel).Error; err != nil {
				return err
			}
		}
	} else {
		if err := tx.Create(userModel).Error; err != nil {
			return err
		}
		userInfo.ID = userModel.ID
		userInfo.CreatedAt = userModel.CreatedAt.Time
	}

	return nil
}

func (repo *userRepository) userFromModel(u *user.User) *model.User {
	telephone, _ := sharedVO.NewTelephone(u.Telephone)
	wxUnionId, _ := valueobject.NewWxUnionID(u.WxUnionID)
	wxMiniProgramOpenId, _ := valueobject.NewWxMiniProgramOpenID(u.WxMiniProgramOpenID)
	wxOfficialOpenId, _ := valueobject.NewWxOfficialOpenID(u.WxOfficialOpenID)
	avatar, _ := valueobject.NewAvatar(u.Avatar)
	return &model.User{
		ID:                  u.ID,
		Telephone:           telephone,
		WxUnionID:           wxUnionId,
		WxMiniProgramOpenID: wxMiniProgramOpenId,
		WxOfficialOpenID:    wxOfficialOpenId,
		Nickname:            u.Nickname,
		Avatar:              avatar,
		CreatedAt:           u.CreatedAt.Time,
	}
}

func (repo *userRepository) userListFromModel(userList []user.User) []model.User {
	result := make([]model.User, len(userList))
	for i, u := range userList {
		result[i] = *repo.userFromModel(&u)
	}
	return result
}

func (repo *userRepository) toModel(u *model.User) *user.User {
	return &user.User{
		DistributedPrimaryKeyID: sharedModel.DistributedPrimaryKeyID{ID: u.ID},
		Telephone:               u.Telephone.String(),
		WxUnionID:               u.WxUnionID.String(),
		WxMiniProgramOpenID:     u.WxMiniProgramOpenID.String(),
		WxOfficialOpenID:        u.WxOfficialOpenID.String(),
		Nickname:                u.Nickname,
		Avatar:                  u.Avatar.RelativePath(),
		Status:                  u.Status.Uint(),
	}
}
