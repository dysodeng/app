package transactions

import (
	"context"

	"gorm.io/gorm"
)

// TxKey 事务上下文的key
type TxKey struct{}

// TransactionManager 事务管理器
type TransactionManager interface {
	// Transaction 开启事务上下文
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
	// GetTx 从上下文中获取事务
	GetTx(ctx context.Context) *gorm.DB
}
