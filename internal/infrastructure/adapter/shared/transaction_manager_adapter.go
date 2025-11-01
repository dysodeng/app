package shared

import (
	"context"

	domainPort "github.com/dysodeng/app/internal/domain/shared/port"
	infraTx "github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
)

// TransactionManagerAdapter 事务管理器适配器
type TransactionManagerAdapter struct {
	tx infraTx.TransactionManager
}

func NewTransactionManagerAdapter(tx infraTx.TransactionManager) domainPort.TransactionManager {
	return &TransactionManagerAdapter{tx: tx}
}

func (a *TransactionManagerAdapter) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return a.tx.Transaction(ctx, fn)
}
