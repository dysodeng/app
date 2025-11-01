package port

import "context"

// TransactionManager 事务管理端口
type TransactionManager interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
