package health

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
)

// healthServer 容器环境健康检查服务
type healthServer struct {
	listener net.Listener
}

func NewServer() server.Server {
	return &healthServer{}
}

func (hs *healthServer) IsEnabled() bool {
	return config.Server.Health.Enabled
}

func (hs *healthServer) Serve() {
	log.Println("start health server...")

	var err error
	hs.listener, err = net.Listen("tcp4", fmt.Sprintf(":%s", config.Server.Health.Port))
	if err != nil {
		panic(err)
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	log.Printf("listening health service on %s\n", ":"+config.Server.Health.Port)

	go func() {
		for {
			conn, err := hs.listener.Accept()
			if err != nil {
				break
			}
			go health(conn)
		}
	}()
}

func (hs *healthServer) Shutdown() {
	log.Println("shutdown health server...")
	if hs.listener != nil {
		_ = hs.listener.Close()
	}
}

func health(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	_, _ = conn.Write([]byte("running"))
}
