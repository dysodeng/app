package di

import (
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/shared/db"
	"github.com/dysodeng/app/internal/interfaces/grpc"
	"github.com/dysodeng/app/internal/interfaces/http"
	"github.com/dysodeng/app/internal/interfaces/http/handler"
	"github.com/dysodeng/app/internal/interfaces/websocket"
)

// App 应用程序
type App struct {
	Config         *config.Config
	TxManager      transactions.TransactionManager
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
	moduleRegistry *ModuleRegistrar,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wsServer *websocket.Server,
	eventBus *event.Bus,
) *App {
	return &App{
		Config:         config,
		TxManager:      txManager,
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
func ProvideDB(cfg *config.Config) (transactions.TransactionManager, error) {
	tx, err := db.Initialize(cfg)
	if err != nil {
		return nil, err
	}
	return transactions.NewGormTransactionManager(tx), nil
}

// ProvideHTTPServer 提供HTTP服务器
func ProvideHTTPServer(config *config.Config, moduleRegistry *ModuleRegistrar) *http.Server {
	var handlers []handler.Handler
	for _, module := range moduleRegistry.GetAllModules() {
		handlers = append(handlers, module.Handlers()...)
	}
	return http.NewServer(config, handlers)
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
