package provider

import (
	"context"
	"log/slog"

	"go.uber.org/zap"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/migration"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/pkg/db"
	"github.com/dysodeng/app/internal/infrastructure/pkg/logger"
	"github.com/dysodeng/app/internal/infrastructure/pkg/mq"
	"github.com/dysodeng/app/internal/infrastructure/pkg/redis"
	"github.com/dysodeng/app/internal/infrastructure/pkg/storage"
	"github.com/dysodeng/app/internal/infrastructure/pkg/telemetry"
)

// ProvideConfig 提供配置
func ProvideConfig() (*config.LoadResult, error) {
	result, err := config.Load("configs/config.yaml")
	if err != nil {
		return nil, err
	}
	if result.Source != "local" {
		slog.Info("配置已从配置中心加载",
			"source", result.Source,
			"path", result.SourcePath)
	}
	return result, nil
}

// ProvideConfigValue 提供当前加载的配置值。
func ProvideConfigValue(result *config.LoadResult) *config.Config {
	return result.Config
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
func ProvideMessageQueue(cfg *config.Config, _ *telemetry.Monitor) (*mq.Client, error) {
	return mq.Init(cfg)
}

// ProvideStorage 提供文件存储
func ProvideStorage(cfg *config.Config) (*storage.Storage, error) {
	return storage.Init(cfg)
}

// ProvideEventBus 提供事件总线
func ProvideEventBus(mq *mq.Client) event.Bus {
	return event.NewMQEventBus(mq)
}
