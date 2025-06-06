package di

import (
	"github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
	"github.com/dysodeng/app/internal/infrastructure/persistence/decorator"
	userRepository "github.com/dysodeng/app/internal/infrastructure/persistence/repository/user"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
)

// ProvideUserRepository 提供用户仓储装饰器
func ProvideUserRepository(
	txManager transactions.TransactionManager,
	cache contract.TypedCache,
) repository.UserRepository {
	// 直接创建原始仓储
	baseRepo := userRepository.NewUserRepository(txManager)
	// 带缓存的用户仓储（缓存仓储装饰器）
	return decorator.NewUserCacheDecorator(baseRepo, cache)
}
