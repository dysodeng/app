package di

import (
	"context"

	"github.com/dysodeng/mq/contract"
	"go.uber.org/zap"

	diEvent "github.com/dysodeng/app/internal/di/event"
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	eventServer "github.com/dysodeng/app/internal/infrastructure/server/event"
	"github.com/dysodeng/app/internal/infrastructure/server/grpc"
	"github.com/dysodeng/app/internal/infrastructure/server/http"
	"github.com/dysodeng/app/internal/infrastructure/server/websocket"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/infrastructure/shared/errors"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry"
	GRPC "github.com/dysodeng/app/internal/interfaces/grpc"
	HTTP "github.com/dysodeng/app/internal/interfaces/http"
	webSocket "github.com/dysodeng/app/internal/interfaces/websocket"
)

// App 应用程序
type App struct {
	Config               *config.Config
	Monitor              *telemetry.Monitor
	Logger               *zap.Logger
	TxManager            transactions.TransactionManager
	RedisClient          redis.Client
	MessageQueue         contract.MQ
	Storage              *storage.Storage
	HandlerRegistry      *HTTP.HandlerRegistry
	WebSocketRegistry    *webSocket.WebSocket
	EventHandlerRegistry *diEvent.HandlerRegistry
	GRPCServiceRegistry  *GRPC.ServiceRegistry
	HTTPServer           *http.Server
	GRPCServer           *grpc.Server
	WSServer             *websocket.Server
	EventBus             event.Bus
	EventConsumer        *event.ConsumerService
	EventServer          *eventServer.Server
}

// NewApp 创建应用程序
func NewApp(
	config *config.Config,
	monitor *telemetry.Monitor,
	logger *zap.Logger,
	txManager transactions.TransactionManager,
	redisClient redis.Client,
	messageQueue contract.MQ,
	storage *storage.Storage,
	handlerRegistry *HTTP.HandlerRegistry,
	webSocketRegistry *webSocket.WebSocket,
	eventHandlerRegistry *diEvent.HandlerRegistry,
	gRPCServiceRegistry *GRPC.ServiceRegistry,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wsServer *websocket.Server,
	eventBus event.Bus,
	eventConsumer *event.ConsumerService,
	eventServer *eventServer.Server,
) *App {
	return &App{
		Config:               config,
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
		EventBus:             eventBus,
		EventConsumer:        eventConsumer,
		EventServer:          eventServer,
	}
}

// Stop 停止应用相关服务
func (app *App) Stop(ctx context.Context) error {
	return errors.NewPipelineWithContext(ctx).Then(db.Close).Then(func() error {
		return app.RedisClient.Close()
	}).Then(func() error {
		return app.MessageQueue.Close()
	}).ExecuteParallel()
}
