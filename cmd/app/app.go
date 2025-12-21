package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dysodeng/app/internal/di"
	"github.com/dysodeng/app/internal/infrastructure/server"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
)

type app struct {
	ctx     context.Context
	mainApp *di.App
	servers []server.Server
}

func newApp(ctx context.Context) *app {
	return &app{
		ctx: ctx,
	}
}

func (app *app) run() {
	// 应用初始化
	app.initialize()

	// 启动服务
	app.serve()

	// 等待中断信息并优雅地关闭服务
	app.waitForInterruptSignal()
}

func (app *app) initialize() {
	mainApp, err := di.InitApp(app.ctx)
	if err != nil {
		logger.Fatal(app.ctx, "应用初始化失败", logger.ErrorField(err))
	}
	app.mainApp = mainApp
}

func (app *app) registerServer(servers ...server.Server) {
	for _, svc := range servers {
		if svc.IsEnabled() {
			app.servers = append(app.servers, svc)
		}
	}
}

func (app *app) serve() {
	logger.Info(app.ctx, "start app server...")

	// 注册服务
	app.registerServer(
		app.mainApp.EventServer,
		app.mainApp.GRPCServer,
		app.mainApp.HTTPServer,
		app.mainApp.WSServer,
		app.mainApp.HealthServer,
	)

	// 启动服务
	for _, serverIns := range app.servers {
		if err := serverIns.Start(); err != nil {
			logger.Fatal(app.ctx, fmt.Sprintf("%s服务启动失败", serverIns.Name()), logger.ErrorField(err))
		}
		logger.Info(app.ctx, fmt.Sprintf("%s服务已启动", serverIns.Name()), logger.AddField("addr", serverIns.Addr()))
	}
}

func (app *app) waitForInterruptSignal() {
	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(app.ctx, "正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务
	for _, serverIns := range app.servers {
		if err := serverIns.Stop(ctx); err != nil {
			logger.Error(ctx, fmt.Sprintf("%s服务关闭失败", serverIns.Name()), logger.ErrorField(err))
		}
		logger.Info(ctx, fmt.Sprintf("%s服务已关闭", serverIns.Name()))
	}

	if err := app.mainApp.Stop(ctx); err != nil {
		logger.Error(ctx, "应用停止失败", logger.ErrorField(err))
	}
}

func Execute() {
	ctx := context.Background()
	newApp(ctx).run()
}
