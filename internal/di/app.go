package di

import (
	"context"

	"go.uber.org/zap"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/server/grpc"
	"github.com/dysodeng/app/internal/infrastructure/server/http"
	"github.com/dysodeng/app/internal/infrastructure/server/websocket"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/errors"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry"
)

// App 应用程序
type App struct {
	Config         *config.Config
	Monitor        *telemetry.Monitor
	Logger         *zap.Logger
	TxManager      transactions.TransactionManager
	RedisClient    redis.Client
	Storage        *storage.Storage
	ModuleRegistry *ModuleRegistrar
	HTTPServer     *http.Server
	GRPCServer     *grpc.Server
	WSServer       *websocket.Server
	EventBus       *event.Bus
}

// NewApp 创建应用程序
func NewApp(
	config *config.Config,
	monitor *telemetry.Monitor,
	logger *zap.Logger,
	txManager transactions.TransactionManager,
	redisClient redis.Client,
	storage *storage.Storage,
	moduleRegistry *ModuleRegistrar,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wsServer *websocket.Server,
	eventBus *event.Bus,
) *App {
	return &App{
		Config:         config,
		Monitor:        monitor,
		Logger:         logger,
		TxManager:      txManager,
		RedisClient:    redisClient,
		Storage:        storage,
		ModuleRegistry: moduleRegistry,
		HTTPServer:     httpServer,
		GRPCServer:     grpcServer,
		WSServer:       wsServer,
		EventBus:       eventBus,
	}
}

// Stop 停止应用相关服务
func (app *App) Stop(ctx context.Context) error {
	pipeline := errors.NewPipelineWithContext(ctx)
	return pipeline.Then(func() error {
		return db.Close()
	}).Then(func() error {
		return app.RedisClient.Close()
	}).ExecuteParallel()
}
