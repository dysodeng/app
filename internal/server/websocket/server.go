package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dysodeng/app/internal/api/websocket"
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/ws"
	"github.com/dysodeng/app/internal/server"
	"github.com/pkg/errors"
)

// websocketServer websocket服务
type websocketServer struct {
	wss *http.Server
}

func NewServer() server.Interface {
	return &websocketServer{}
}

func (server *websocketServer) Serve() {
	if !config.Server.Websocket.Enabled {
		return
	}
	log.Println("start websocket server...")

	// websocket 客户端连接hub
	ws.HubBus = ws.NewHub()
	go ws.HubBus.Run()
	ws.HubBus.SetMessageHandler(&websocket.MessageHandler{})

	mux := http.NewServeMux()
	mux.HandleFunc("/ws/v1/message", func(w http.ResponseWriter, r *http.Request) {
		ws.WebsocketServe(w, r)
	})
	mux.HandleFunc("/ws/v1/metrics", func(w http.ResponseWriter, r *http.Request) {
		ws.Metrics(w)
	})

	server.wss = &http.Server{
		Addr:              ":" + config.Server.Websocket.Port,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		if err := server.wss.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("websocket server error: %s\n", err)
		}
	}()

	log.Printf("websocket service listening and serving 0.0.0.0:%s\n", config.Server.Websocket.Port)
}

func (server *websocketServer) Shutdown() {
	if !config.Server.Websocket.Enabled {
		return
	}
	log.Println("shutdown websocket server...")

	err := server.wss.Shutdown(context.Background())
	if err != nil {
		log.Printf("websocket server shutdown fiald:%s", err)
	}
}
