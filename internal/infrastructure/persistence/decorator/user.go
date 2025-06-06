package decorator

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dysodeng/app/internal/domain/user/model"
	"github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/strategy"
)

// UserCacheDecorator 用户仓储装饰器
type UserCacheDecorator struct {
	repository repository.UserRepository // 原始仓储
	cache      contract.TypedCache
	strategy   contract.CacheStrategy
}

func NewUserCacheDecorator(
	repo repository.UserRepository,
	cache contract.TypedCache,
) repository.UserRepository {
	return &UserCacheDecorator{
		repository: repo,
		cache:      cache,
		strategy:   strategy.NewUserStrategy(),
	}
}

func (d *UserCacheDecorator) Info(ctx context.Context, id uint64) (*model.User, error) {
	// 缓存逻辑
	cacheKey := d.strategy.GetCacheKey("user:info", fmt.Sprintf("%d", id))

	var user *model.User
	if err := d.cache.GetObject(ctx, cacheKey, &user); err == nil {
		return user, nil
	}

	// 调用原始仓储
	user, err := d.repository.Info(ctx, id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	_ = d.cache.PutObject(ctx, cacheKey, user, time.Hour)

	return user, nil
}

func (d *UserCacheDecorator) ListUser(ctx context.Context, query repository.UserListQuery) ([]model.User, int64, error) {
	// 生成查询条件的哈希作为缓存键的一部分
	queryHash := d.generateQueryHash(query)
	cacheKey := d.strategy.GetCacheKey("user:list", queryHash)

	// 定义缓存结构
	type CachedUserList struct {
		Users []model.User `json:"users"`
		Total int64        `json:"total"`
	}

	// 尝试从缓存获取
	var cached CachedUserList
	if err := d.cache.GetObject(ctx, cacheKey, &cached); err == nil {
		return cached.Users, cached.Total, nil
	}

	// 缓存未命中，调用原始仓储
	users, total, err := d.repository.ListUser(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 将结果写入缓存（列表查询缓存时间可以设置短一些）
	cachedData := CachedUserList{
		Users: users,
		Total: total,
	}
	_ = d.cache.PutObject(ctx, cacheKey, cachedData, 10*time.Minute)

	return users, total, nil
}

// generateQueryHash 生成查询条件的哈希值
func (d *UserCacheDecorator) generateQueryHash(query repository.UserListQuery) string {
	// 使用 fmt.Sprintf 生成查询条件的字符串表示
	queryStr := fmt.Sprintf("%s_%s_%s_%d_%s_%s_%d_%d",
		query.Telephone,
		query.RealName,
		query.Nickname,
		query.Status,
		query.OrderBy,
		query.OrderType,
		query.Page,
		query.PageSize,
	)

	// 生成 MD5 哈希
	hash := md5.Sum([]byte(queryStr))
	return fmt.Sprintf("%x", hash)
}

func (d *UserCacheDecorator) CreateUser(ctx context.Context, userInfo *model.User) error {
	err := d.repository.CreateUser(ctx, userInfo)
	if err != nil {
		return err
	}

	// 清理列表缓存
	d.invalidateListCache(ctx)
	return nil
}

func (d *UserCacheDecorator) UpdateUser(ctx context.Context, userInfo *model.User) error {
	err := d.repository.UpdateUser(ctx, userInfo)
	if err != nil {
		return err
	}

	// 删除用户信息缓存
	cacheKey := d.strategy.GetCacheKey("user:info", fmt.Sprintf("%d", userInfo.ID))
	_ = d.cache.DeleteObject(ctx, cacheKey)

	// 清理列表缓存
	d.invalidateListCache(ctx)

	return nil
}

func (d *UserCacheDecorator) DeleteUser(ctx context.Context, userId uint64) error {
	if err := d.repository.DeleteUser(ctx, userId); err != nil {
		return err
	}

	// 删除用户信息缓存
	cacheKey := d.strategy.GetCacheKey("user:info", fmt.Sprintf("%d", userId))
	_ = d.cache.DeleteObject(ctx, cacheKey)

	// 清理列表缓存
	d.invalidateListCache(ctx)

	return nil
}

// invalidateListCache 清理列表缓存
func (d *UserCacheDecorator) invalidateListCache(ctx context.Context) {
	// 使用批量删除清理所有用户列表缓存
	listCachePrefix := d.strategy.GetCacheKey("user:list", "")
	// 注意：这里需要确保你的 cache 实现支持按前缀批量删除
	// 如果不支持，可以考虑使用版本号机制或者设置较短的过期时间
	if baseCache, ok := d.cache.(interface {
		BatchDelete(ctx context.Context, prefix string) error
	}); ok {
		_ = baseCache.BatchDelete(ctx, listCachePrefix)
	}
}
