package transactions

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/db"
	"gorm.io/gorm"
)

type gormTransactionManager struct {
	db *gorm.DB
}

// NewGormTransactionManager 创建GORM事务管理器
func NewGormTransactionManager() TransactionManager {
	return &gormTransactionManager{
		db: db.DB(),
	}
}

// Transaction 开启事务上下文
func (tm *gormTransactionManager) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	// 检查是否已经在事务中
	if existingTx, ok := ctx.Value(TxKey{}).(*gorm.DB); ok && existingTx != nil {
		return fn(ctx)
	}

	// 开启新事务
	err := tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建新的事务上下文
		txCtx := context.WithValue(ctx, TxKey{}, tx)

		// 执行事务函数
		if err := fn(txCtx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// GetTx 从上下文中获取事务
func (tm *gormTransactionManager) GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TxKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return tm.db.WithContext(ctx)
}
