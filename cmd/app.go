package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/redis"
	"github.com/dysodeng/app/internal/pkg/telemetry/metrics"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/app/internal/server/cron"
	"github.com/dysodeng/app/internal/server/grpc"
	"github.com/dysodeng/app/internal/server/health"
	"github.com/dysodeng/app/internal/server/http"
	"github.com/dysodeng/app/internal/server/task"
	"github.com/dysodeng/app/internal/server/websocket"
	"github.com/dysodeng/app/migrations/migration"
)

type app struct {
	serverList []server.Interface
}

func newApp() *app {
	return &app{}
}

func (app *app) run() {
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

func (app *app) config() {
	config.Load()
}

func (app *app) monitor() {
	if err := trace.Init(); err != nil {
		panic(err)
	}
	if err := metrics.Init(); err != nil {
		panic(err)
	}
}

func (app *app) initialize() {
	db.Initialize()
	redis.Initialize()

	// db migration
	if config.Database.Migration.Enabled {
		if err := migration.Migrate(db.DB()); err != nil {
			panic(err)
		}
	}
}

func (app *app) start() {
	log.Println("start app server...")

	// 注册服务
	app.registerServer(
		task.NewServer(),
		cron.NewServer(),
		grpc.NewServer(),
		http.NewServer(),
		websocket.NewServer(),
		health.NewServer(),
	)

	// 启动服务
	for _, serverIns := range app.serverList {
		serverIns.Serve()
	}
}

func (app *app) registerServer(servers ...server.Interface) {
	for _, s := range servers {
		if s.IsEnabled() {
			app.serverList = append(app.serverList, s)
		}
	}
}

func (app *app) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	for _, serverIns := range app.serverList {
		serverIns.Shutdown()
	}

	db.Close()
	redis.Close()

	log.Println("shutdown app server...")
}

func Execute() {
	newApp().run()
}
