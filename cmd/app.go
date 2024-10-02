package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/dysodeng/app/internal/service/provider"

	"github.com/dysodeng/app/migrations/migration"

	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/redis"
	"github.com/dysodeng/app/internal/server/cron"
	"github.com/dysodeng/app/internal/server/health"
	"github.com/dysodeng/app/internal/server/http"
	"github.com/dysodeng/app/internal/server/task"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
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

func (app *App) initialize() {
	db.Initialize()
	redis.Initialize()

	// db migration
	if config.Database.Migration.Enabled {
		if err := migration.Migrate(db.DB()); err != nil {
			panic(err)
		}
	}

	// service container
	provider.ServiceProvider()
}

func (app *App) start() {
	log.Println("start app server...")

	// 注册服务
	app.registerServer(
		task.NewServer(),
		cron.NewServer(),
		http.NewServer(),
		health.NewServer(),
	)

	// 启动服务
	for _, serverIns := range app.serverList {
		serverIns.Serve()
	}
}

func (app *App) registerServer(servers ...server.Interface) {
	app.serverList = append(app.serverList, servers...)
}

func (app *App) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	app.shutdown()
}

func (app *App) shutdown() {
	for _, serverIns := range app.serverList {
		serverIns.Shutdown()
	}
	log.Println("shutdown app server...")
}
