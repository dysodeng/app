package cache

import (
	"context"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/google/uuid"

	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
	"github.com/dysodeng/app/internal/domain/user/model"
	userDomainRepo "github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/domain/user/valueobject"
	"github.com/dysodeng/app/internal/infrastructure/config"
	persistCache "github.com/dysodeng/app/internal/infrastructure/persistence/cache"
	userRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/user"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	sharedModel "github.com/dysodeng/app/internal/infrastructure/shared/model"
)

// 缓存DTO，避免领域值对象的私有字段导致JSON不完整
type userCacheDTO struct {
	ID                  uuid.UUID `json:"id"`
	Telephone           string    `json:"telephone"`
	WxUnionID           string    `json:"wx_union_id"`
	WxMiniProgramOpenID string    `json:"wx_mini_program_openid"`
	WxOfficialOpenID    string    `json:"wx_official_openid"`
	Nickname            string    `json:"nickname"`
	Avatar              string    `json:"avatar"`
	Status              uint8     `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

type cachedUserRepository struct {
	next     userDomainRepo.UserRepository
	cache    *persistCache.TypedCache[userCacheDTO]
	cacheTTL time.Duration
}

func NewCachedUserRepository(txManager transactions.TransactionManager) userDomainRepo.UserRepository {
	driver := config.GlobalConfig.Cache.Driver
	cacheTTL := 10 * time.Minute
	tc := persistCache.NewTypedCacheWith[userCacheDTO](driver, "user", cacheTTL, true)
	return &cachedUserRepository{
		next:     userRepository.NewUserRepository(txManager),
		cache:    tc,
		cacheTTL: cacheTTL,
	}
}

func toDTO(u *model.User) *userCacheDTO {
	return &userCacheDTO{
		ID:                  u.ID,
		Telephone:           u.Telephone.Value(),
		WxUnionID:           u.WxUnionID.Value(),
		WxMiniProgramOpenID: u.WxMiniProgramOpenID.Value(),
		WxOfficialOpenID:    u.WxOfficialOpenID.Value(),
		Nickname:            u.Nickname,
		Avatar:              u.Avatar.FullURL(),
		Status:              u.Status.Uint(),
		CreatedAt:           u.CreatedAt,
	}
}

func toDomain(dto *userCacheDTO) *model.User {
	tel, _ := sharedVO.NewTelephone(dto.Telephone)
	union, _ := valueobject.NewWxUnionID(dto.WxUnionID)
	wxmp, _ := valueobject.NewWxMiniProgramOpenID(dto.WxMiniProgramOpenID)
	wxoff, _ := valueobject.NewWxOfficialOpenID(dto.WxOfficialOpenID)
	avatar, _ := valueobject.NewAvatar(dto.Avatar)
	return &model.User{
		ID:                  dto.ID,
		Telephone:           tel,
		WxUnionID:           union,
		WxMiniProgramOpenID: wxmp,
		WxOfficialOpenID:    wxoff,
		Nickname:            dto.Nickname,
		Avatar:              avatar,
		Status:              sharedModel.BinaryStatusByUint(dto.Status),
		CreatedAt:           dto.CreatedAt,
	}
}

func (r *cachedUserRepository) tagsFor(u *model.User) []string {
	var tags []string
	if u.ID != uuid.Nil {
		tags = append(tags, "user:"+u.ID.String())
	}
	if v := u.Telephone.Value(); v != "" {
		tags = append(tags, "tel:"+v)
	}
	if v := u.WxUnionID.Value(); v != "" {
		tags = append(tags, "union:"+v)
	}
	if v := u.WxMiniProgramOpenID.Value(); v != "" {
		tags = append(tags, "mp:"+v)
	}
	if v := u.WxOfficialOpenID.Value(); v != "" {
		tags = append(tags, "off:"+v)
	}
	return tags
}

func (r *cachedUserRepository) invalidateByUser(ctx context.Context, u *model.User) {
	if u == nil {
		return
	}
	_ = r.cache.InvalidateTags(ctx, r.tagsFor(u)...)
}

func (r *cachedUserRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	base := "id:" + id.String()
	// 命中缓存
	if dto, ok, _ := r.cache.Get(ctx, base, "user:"+id.String()); ok && dto.ID != uuid.Nil {
		return toDomain(&dto), nil
	}
	// 加载并写入缓存（仅缓存有效实体）
	u, err := r.next.FindById(ctx, id)
	if err != nil || u == nil || u.ID == uuid.Nil {
		return u, err
	}
	_ = r.cache.Set(ctx, base, *toDTO(u), r.cacheTTL, "user:"+u.ID.String())
	return u, nil
}

func (r *cachedUserRepository) FindByTelephone(ctx context.Context, telephone string) (*model.User, error) {
	base := "tel:" + telephone
	if dto, ok, _ := r.cache.Get(ctx, base, "tel:"+telephone); ok && dto.ID != uuid.Nil {
		return toDomain(&dto), nil
	}
	u, err := r.next.FindByTelephone(ctx, telephone)
	if err != nil || u == nil || u.ID == uuid.Nil {
		return u, err
	}
	_ = r.cache.Set(ctx, base, *toDTO(u), r.cacheTTL, "tel:"+telephone)
	return u, nil
}

func (r *cachedUserRepository) FindByUnionId(ctx context.Context, unionId string) (*model.User, error) {
	base := "union:" + unionId
	if dto, ok, _ := r.cache.Get(ctx, base, "union:"+unionId); ok && dto.ID != uuid.Nil {
		return toDomain(&dto), nil
	}
	u, err := r.next.FindByUnionId(ctx, unionId)
	if err != nil || u == nil || u.ID == uuid.Nil {
		return u, err
	}
	_ = r.cache.Set(ctx, base, *toDTO(u), r.cacheTTL, "union:"+unionId)
	return u, nil
}

func (r *cachedUserRepository) FindByOpenId(ctx context.Context, platform, openId string) (*model.User, error) {
	tag := "openid:" + platform + ":" + openId
	base := tag
	if dto, ok, _ := r.cache.Get(ctx, base, tag); ok && dto.ID != uuid.Nil {
		return toDomain(&dto), nil
	}
	u, err := r.next.FindByOpenId(ctx, platform, openId)
	if err != nil || u == nil || u.ID == uuid.Nil {
		return u, err
	}
	err = r.cache.Set(ctx, base, *toDTO(u), r.cacheTTL, tag)
	if err != nil {
		logger.Error(ctx, "缓存失败", logger.ErrorField(err))
	}
	return u, nil
}

func (r *cachedUserRepository) Save(ctx context.Context, userInfo *model.User) error {
	if err := r.next.Save(ctx, userInfo); err != nil {
		return err
	}
	// 标签版本失效，避免扫描删除
	r.invalidateByUser(ctx, userInfo)
	return nil
}
