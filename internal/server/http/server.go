package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dysodeng/app/internal/api/http/router"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
	"github.com/gin-gonic/gin"
)

// httpServer http服务
type httpServer struct {
	server *http.Server
}

func NewServer() server.Interface {
	return &httpServer{}
}

func (httpServer *httpServer) init() {
	// env mode
	switch config.App.Env {
	case config.Dev:
		gin.SetMode(gin.DebugMode)
	case config.Test:
		gin.SetMode(gin.TestMode)
	case config.Prod:
		gin.SetMode(gin.ReleaseMode)
	}

	// error logger
	logFilename := config.LogPath + "/gin.error.log"
	errLogFile, _ := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	gin.DefaultErrorWriter = io.MultiWriter(errLogFile, os.Stderr)

	httpServer.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", config.Server.Http.Port),
		Handler: router.Router(),
	}
}

func (httpServer *httpServer) Serve() {
	if !config.Server.Http.Enabled {
		return
	}
	log.Println("start http server...")

	httpServer.init()

	defer func() {
		if err := recover(); err != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err = httpServer.server.Shutdown(ctx); err != nil {
				log.Fatalf("http server shutdown fiald: %+v\n", err)
			}
		}
	}()

	log.Printf("http service listening and serving 0.0.0.0:%s\n", config.Server.Http.Port)

	go func() {
		if err := httpServer.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server start fiald: %s\n", err)
		}
	}()
}

func (httpServer *httpServer) Shutdown() {
	if !config.Server.Http.Enabled {
		return
	}
	log.Println("shutdown http server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.server.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown fiald:%s", err)
	}
}
