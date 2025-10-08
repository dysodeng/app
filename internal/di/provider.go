package di

import (
	"context"

	"github.com/dysodeng/mq/contract"
	"go.uber.org/zap"

	diEvent "github.com/dysodeng/app/internal/di/event"
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/migration"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	eventServer "github.com/dysodeng/app/internal/infrastructure/server/event"
	"github.com/dysodeng/app/internal/infrastructure/server/grpc"
	"github.com/dysodeng/app/internal/infrastructure/server/health"
	"github.com/dysodeng/app/internal/infrastructure/server/http"
	"github.com/dysodeng/app/internal/infrastructure/server/websocket"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/mq"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry"
	GRPC "github.com/dysodeng/app/internal/interfaces/grpc"
	HTTP "github.com/dysodeng/app/internal/interfaces/http"
	webSocket "github.com/dysodeng/app/internal/interfaces/websocket"
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

// ProvideHTTPServer 提供HTTP服务器
func ProvideHTTPServer(cfg *config.Config, handlerRegistry *HTTP.HandlerRegistry) *http.Server {
	return http.NewServer(cfg, handlerRegistry)
}

// ProvideGRPCServer 提供gRPC服务器
func ProvideGRPCServer(ctx context.Context, cfg *config.Config, serviceRegistry *GRPC.ServiceRegistry) *grpc.Server {
	return grpc.NewServer(ctx, cfg, serviceRegistry)
}

// ProvideWebSocketServer 提供WebSocket服务器
func ProvideWebSocketServer(cfg *config.Config, ws *webSocket.WebSocket) *websocket.Server {
	return websocket.NewServer(cfg, ws)
}

// ProvideHealthServer 提供容器环境健康检查服务
func ProvideHealthServer(cfg *config.Config) *health.Server {
	return health.NewServer(cfg)
}

// ProvideTypedEventBus 提供类型化事件总线
func ProvideTypedEventBus(mq contract.MQ) event.Bus {
	return event.NewMQEventBus(mq.Producer())
}

// ProvideEventConsumerService 提供事件消费者服务
func ProvideEventConsumerService(mq contract.MQ, logger *zap.Logger) *event.ConsumerService {
	return event.NewEventConsumerService(mq.Consumer(), logger)
}

// ProvideEventServer 提供Event服务器
func ProvideEventServer(
	cfg *config.Config,
	eventConsumer *event.ConsumerService,
	registry *diEvent.HandlerRegistry,
) *eventServer.Server {
	return eventServer.NewEventServer(cfg, eventConsumer, registry)
}
