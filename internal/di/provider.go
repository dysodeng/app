package di

import (
	"context"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/migration"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/server/grpc"
	"github.com/dysodeng/app/internal/infrastructure/server/http"
	"github.com/dysodeng/app/internal/infrastructure/server/websocket"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
	"github.com/dysodeng/app/internal/interfaces/http/handler"
)

// App 应用程序
type App struct {
	Config         *config.Config
	TxManager      transactions.TransactionManager
	RedisClient    redis.Client
	ModuleRegistry *ModuleRegistrar
	HTTPServer     *http.Server
	GRPCServer     *grpc.Server
	WSServer       *websocket.Server
	EventBus       *event.Bus
}

// NewApp 创建应用程序
func NewApp(
	config *config.Config,
	txManager transactions.TransactionManager,
	redisClient redis.Client,
	moduleRegistry *ModuleRegistrar,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wsServer *websocket.Server,
	eventBus *event.Bus,
) *App {
	return &App{
		Config:         config,
		TxManager:      txManager,
		RedisClient:    redisClient,
		ModuleRegistry: moduleRegistry,
		HTTPServer:     httpServer,
		GRPCServer:     grpcServer,
		WSServer:       wsServer,
		EventBus:       eventBus,
	}
}

// ProvideConfig 提供配置
func ProvideConfig() (*config.Config, error) {
	return config.LoadConfig("configs/config.yaml")
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

// ProvideHTTPServer 提供HTTP服务器
func ProvideHTTPServer(config *config.Config, moduleRegistry *ModuleRegistrar) *http.Server {
	var handlers []handler.Handler
	for _, module := range moduleRegistry.GetAllModules() {
		handlers = append(handlers, module.Handlers()...)
	}
	return http.NewServer(config, handlers...)
}

// ProvideGRPCServer 提供gRPC服务器
func ProvideGRPCServer(config *config.Config, moduleRegistry *ModuleRegistrar) *grpc.Server {
	server := grpc.NewServer(config)
	// 注册gRPC服务在这里实现
	return server
}

// ProvideWebSocketServer 提供WebSocket服务器
func ProvideWebSocketServer(config *config.Config) *websocket.Server {
	return websocket.NewServer(config)
}

// ProvideEventBus 提供事件总线
func ProvideEventBus(moduleRegistry *ModuleRegistrar) *event.Bus {
	eventBus := event.NewEventBus()
	for _, module := range moduleRegistry.GetAllModules() {
		for _, eventHandler := range module.EventHandlers() {
			eventBus.RegisterHandler(eventHandler)
		}
	}
	return eventBus
}
