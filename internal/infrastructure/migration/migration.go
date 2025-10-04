package migration

import (
	"context"

	"github.com/go-gormigrate/gormigrate/v2"

	"github.com/dysodeng/app/internal/infrastructure/persistence/entity"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
)

// 定义数据库迁移
var migrations []*gormigrate.Migration

func margeMigrations() {
	migrations = append(migrations, userMigrations...)
	migrations = append(migrations, fileMigrations...)
}

// Migrate 执行数据库迁移
func Migrate(ctx context.Context, tx transactions.TransactionManager) error {
	logger.Info(ctx, "开始数据库迁移")

	margeMigrations()
	if len(migrations) == 0 {
		return nil
	}

	// 自动迁移数据库表结构
	err := gormigrate.New(tx.GetTx(ctx), gormigrate.DefaultOptions, migrations).Migrate()
	if err != nil {
		logger.Error(ctx, "数据库迁移失败", logger.ErrorField(err))
		return err
	}

	logger.Info(ctx, "数据库迁移完成")
	return nil
}

// Rollback 执行数据库回滚
func Rollback(ctx context.Context, tx transactions.TransactionManager, version ...string) error {
	logger.Info(ctx, "开始数据库迁移回滚")

	margeMigrations()
	if len(migrations) == 0 {
		return nil
	}

	var err error
	if len(version) > 0 {
		err = gormigrate.New(tx.GetTx(ctx), gormigrate.DefaultOptions, migrations).RollbackTo(version[0])
	} else {
		err = gormigrate.New(tx.GetTx(ctx), gormigrate.DefaultOptions, migrations).RollbackLast()
	}
	if err != nil {
		logger.Error(ctx, "数据库迁移回滚失败", logger.ErrorField(err))
		return err
	}

	logger.Info(ctx, "数据库迁移回滚完成")
	return nil
}

// Seed 填充初始数据
func Seed(ctx context.Context, tx transactions.TransactionManager) error {
	logger.Info(ctx, "开始填充初始数据")

	// 检查是否已有管理员用户
	var count int64
	tx.GetTx(ctx).Model(&entity.User{}).Count(&count)

	// 如果没有用户，创建一个管理员用户
	if count == 0 {
		adminUser := &entity.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // 密码: password
		}

		if err := tx.GetTx(ctx).Create(adminUser).Error; err != nil {
			logger.Error(ctx, "创建管理员用户失败", logger.ErrorField(err))
			return err
		}

		logger.Info(ctx, "创建管理员用户成功")
	}

	logger.Info(ctx, "初始数据填充完成")
	return nil
}
