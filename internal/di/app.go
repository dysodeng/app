package di

import (
	"context"
	stderrors "errors"

	"go.uber.org/zap"

	diEvent "github.com/dysodeng/app/internal/di/event"
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/pkg/db"
	"github.com/dysodeng/app/internal/infrastructure/pkg/errors"
	appmq "github.com/dysodeng/app/internal/infrastructure/pkg/mq"
	"github.com/dysodeng/app/internal/infrastructure/pkg/redis"
	"github.com/dysodeng/app/internal/infrastructure/pkg/storage"
	"github.com/dysodeng/app/internal/infrastructure/pkg/telemetry"
	eventServer "github.com/dysodeng/app/internal/infrastructure/server/event"
	"github.com/dysodeng/app/internal/infrastructure/server/grpc"
	"github.com/dysodeng/app/internal/infrastructure/server/health"
	"github.com/dysodeng/app/internal/infrastructure/server/http"
	"github.com/dysodeng/app/internal/infrastructure/server/websocket"
	GRPC "github.com/dysodeng/app/internal/interfaces/grpc"
	HTTP "github.com/dysodeng/app/internal/interfaces/http"
	webSocket "github.com/dysodeng/app/internal/interfaces/websocket"
)

// App 应用程序
type App struct {
	Config               *config.Config
	ConfigWatcher        *config.Watcher
	Monitor              *telemetry.Monitor
	Logger               *zap.Logger
	TxManager            transactions.TransactionManager
	RedisClient          redis.Client
	MessageQueue         *appmq.Client
	Storage              *storage.Storage
	HandlerRegistry      *HTTP.HandlerRegistry
	WebSocketRegistry    *webSocket.WebSocket
	EventHandlerRegistry *diEvent.HandlerRegistry
	GRPCServiceRegistry  *GRPC.ServiceRegistry
	HTTPServer           *http.Server
	GRPCServer           *grpc.Server
	WSServer             *websocket.Server
	HealthServer         *health.Server
	EventBus             event.Bus
	EventConsumer        *event.ConsumerService
	EventServer          *eventServer.Server
}

// NewApp 创建应用程序
func NewApp(
	cfg *config.Config,
	loadResult *config.LoadResult,
	monitor *telemetry.Monitor,
	logger *zap.Logger,
	txManager transactions.TransactionManager,
	redisClient redis.Client,
	messageQueue *appmq.Client,
	storage *storage.Storage,
	handlerRegistry *HTTP.HandlerRegistry,
	webSocketRegistry *webSocket.WebSocket,
	eventHandlerRegistry *diEvent.HandlerRegistry,
	gRPCServiceRegistry *GRPC.ServiceRegistry,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wsServer *websocket.Server,
	healthServer *health.Server,
	eventBus event.Bus,
	eventConsumer *event.ConsumerService,
	eventServer *eventServer.Server,
) *App {
	var watcher *config.Watcher
	if loadResult != nil && loadResult.WatchSource != nil {
		watcher = config.NewWatcher(loadResult.WatchSource, loadResult.LocalPath, nil)
	}
	return &App{
		Config:               cfg,
		ConfigWatcher:        watcher,
		Monitor:              monitor,
		Logger:               logger,
		TxManager:            txManager,
		RedisClient:          redisClient,
		MessageQueue:         messageQueue,
		Storage:              storage,
		HandlerRegistry:      handlerRegistry,
		WebSocketRegistry:    webSocketRegistry,
		EventHandlerRegistry: eventHandlerRegistry,
		GRPCServiceRegistry:  gRPCServiceRegistry,
		HTTPServer:           httpServer,
		GRPCServer:           grpcServer,
		WSServer:             wsServer,
		HealthServer:         healthServer,
		EventBus:             eventBus,
		EventConsumer:        eventConsumer,
		EventServer:          eventServer,
	}
}

// StartConfigWatcher 启动配置中心热更新监听。
func (app *App) StartConfigWatcher(ctx context.Context) error {
	if app.ConfigWatcher == nil {
		return nil
	}
	return app.ConfigWatcher.Start(ctx)
}

// Stop 停止应用相关服务
func (app *App) Stop(ctx context.Context) error {
	var watchErr error
	if app.ConfigWatcher != nil {
		watchErr = app.ConfigWatcher.Stop()
	}
	resourceErr := errors.NewPipelineWithContext(ctx).Then(db.Close).Then(func() error {
		return app.RedisClient.Close()
	}).Then(func() error {
		return app.MessageQueue.Close(ctx)
	}).ExecuteParallel()
	return stderrors.Join(watchErr, resourceErr)
}
