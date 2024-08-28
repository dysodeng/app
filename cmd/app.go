package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/app/internal/server/cron"
	"github.com/dysodeng/app/internal/server/health"
	"github.com/dysodeng/app/internal/server/http"
	"github.com/dysodeng/app/internal/server/task"
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

	// 启动服务
	app.registerServer(
		task.NewServer(),
		cron.NewServer(),
		http.NewServer(),
		health.NewServer(),
	)
	app.start()

	// 监听退出信号
	app.keep()
}

func (app *App) config() {
	config.Load()
}

func (app *App) start() {
	log.Println("start app server...")

	// 启动服务
	for _, serverIns := range app.serverList {
		serverIns.Serve()
	}
}

func (app *App) registerServer(servers ...server.Interface) {
	app.serverList = append(app.serverList, servers...)
}

func (app *App) keep() {
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
