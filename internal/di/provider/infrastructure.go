package provider

import (
	"context"

	"github.com/dysodeng/mq/contract"
	"go.uber.org/zap"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/migration"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/mq"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry"
)

// ProvideConfig 提供配置
func ProvideConfig() (*config.Config, error) {
	return config.LoadConfig("configs/config.yaml")
}

// ProvideMonitor 提供可观测性配置
func ProvideMonitor(cfg *config.Config) (*telemetry.Monitor, error) {
	return telemetry.InitMonitor(cfg)
}

// ProvideLogger 提供日志
func ProvideLogger(ctx context.Context, cfg *config.Config) (*zap.Logger, error) {
	logger.InitLogger(cfg.App.Debug)
	logger.Info(ctx, "应用初始化中...")
	return logger.ZapLogger(), nil
}

// ProvideDB 提供数据库
func ProvideDB(ctx context.Context, cfg *config.Config) (transactions.TransactionManager, error) {
	tx, err := db.Initialize(cfg)
	if err != nil {
		return nil, err
	}

	txManager := transactions.NewGormTransactionManager(tx)

	if cfg.Database.Migration.Enabled {
		// 执行数据库迁移
		if err = migration.Migrate(ctx, txManager); err != nil {
			logger.Fatal(ctx, "数据库迁移失败", logger.ErrorField(err))
		}

		// 填充初始数据
		if err = migration.Seed(ctx, txManager); err != nil {
			logger.Fatal(ctx, "初始数据填充失败", logger.ErrorField(err))
		}
	}
	return txManager, nil
}

// ProvideRedis 提供redis
func ProvideRedis(cfg *config.Config) (redis.Client, error) {
	cli, err := redis.Initialize(cfg)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// ProvideMessageQueue 提供消息队列
func ProvideMessageQueue(cfg *config.Config) (contract.MQ, error) {
	return mq.Init(cfg)
}

// ProvideStorage 提供文件存储
func ProvideStorage(cfg *config.Config) (*storage.Storage, error) {
	return storage.Init(cfg)
}

// ProvideEventBus 提供事件总线
func ProvideEventBus(mq contract.MQ) event.Bus {
	return event.NewMQEventBus(mq.Producer())
}
