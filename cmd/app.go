package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/redis"
	telemetryMetrics "github.com/dysodeng/app/internal/pkg/telemetry/metrics"
	telemetryTrace "github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/app/internal/server/cron"
	"github.com/dysodeng/app/internal/server/grpc"
	"github.com/dysodeng/app/internal/server/health"
	"github.com/dysodeng/app/internal/server/http"
	"github.com/dysodeng/app/internal/server/task"
	"github.com/dysodeng/app/internal/server/websocket"
	"github.com/dysodeng/app/migrations/migration"
)

type App struct {
	serverList []server.Interface
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	// 加载配置
	app.config()

	// 开启监控
	app.monitor()

	// 应用初始化
	app.initialize()

	// 启动服务
	app.start()

	// 监听退出信号
	app.listenForShutdown()
}

func (app *App) config() {
	config.Load()
}

func (app *App) monitor() {
	if err := telemetryTrace.Init(); err != nil {
		panic(err)
	}
	if err := telemetryMetrics.Init(); err != nil {
		panic(err)
	}
}

func (app *App) initialize() {
	db.Initialize()
	redis.Initialize()

	// db migration
	if config.Database.Migration.Enabled {
		if err := migration.Migrate(db.DB()); err != nil {
			panic(err)
		}
	}
}

func (app *App) start() {
	log.Println("start app server...")

	// 注册服务
	app.registerServer(
		task.NewServer(),
		cron.NewServer(),
		grpc.NewGrpcServer(),
		http.NewServer(),
		websocket.NewServer(),
		health.NewServer(),
	)

	// 启动服务
	for _, serverIns := range app.serverList {
		serverIns.Serve()
	}
}

func (app *App) registerServer(servers ...server.Interface) {
	for _, s := range servers {
		if s.IsEnabled() {
			app.serverList = append(app.serverList, s)
		}
	}
}

func (app *App) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	for _, serverIns := range app.serverList {
		serverIns.Shutdown()
	}
	log.Println("shutdown app server...")
}
