package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dysodeng/app/internal/di"
	"github.com/dysodeng/app/internal/infrastructure/migration"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
)

func main() {
	ctx := context.Background()

	// 初始化日志
	logger.InitLogger(true)
	logger.Info(ctx, "应用启动中...")

	// 初始化应用
	app, err := di.InitApp()
	if err != nil {
		logger.Fatal(ctx, "初始化应用失败", logger.ErrorField(err))
	}

	// 执行数据库迁移
	if err = migration.Migrate(ctx, app.TxManager); err != nil {
		logger.Fatal(ctx, "数据库迁移失败", logger.ErrorField(err))
	}

	// 填充初始数据
	if err = migration.Seed(ctx, app.TxManager); err != nil {
		logger.Fatal(ctx, "初始数据填充失败", logger.ErrorField(err))
	}

	// 启动HTTP服务
	go func() {
		if err = app.HTTPServer.Start(); err != nil {
			logger.Error(ctx, "HTTP服务启动失败", logger.ErrorField(err))
		}
	}()
	logger.Info(ctx, "HTTP服务已启动", logger.AddField("addr", app.HTTPServer.Addr()))

	// 启动gRPC服务
	go func() {
		if err = app.GRPCServer.Start(); err != nil {
			logger.Error(ctx, "gRPC服务启动失败", logger.ErrorField(err))
		}
	}()
	logger.Info(ctx, "gRPC服务已启动", logger.AddField("addr", app.HTTPServer.Addr()))

	// 启动WebSocket服务
	go func() {
		if err = app.WSServer.Start(); err != nil {
			logger.Error(ctx, "WebSocket服务启动失败", logger.ErrorField(err))
		}
	}()
	logger.Info(ctx, "WebSocket服务已启动", logger.AddField("addr", app.WSServer.Addr()))

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务
	if err = app.HTTPServer.Stop(ctx); err != nil {
		logger.Error(ctx, "HTTP服务关闭失败", logger.ErrorField(err))
	}

	// 关闭gRPC服务
	app.GRPCServer.Stop(ctx)

	// 关闭WebSocket服务
	if err = app.WSServer.Stop(ctx); err != nil {
		logger.Error(ctx, "WebSocket服务关闭失败", logger.ErrorField(err))
	}

	logger.Info(ctx, "服务已关闭")
}
