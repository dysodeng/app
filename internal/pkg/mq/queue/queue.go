package queue

import (
	"github.com/dysodeng/app/internal/pkg/mq"
	"log"
	"time"

	"github.com/dysodeng/mq/contract"
)

var queueProducer contract.Producer

func init() {
	p, err := mq.NewMessageQueueProducer(&contract.Pool{
		MinConn:     2,
		MaxConn:     2,
		MaxIdleConn: 2,
		IdleTimeout: time.Hour,
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	queueProducer = p
}

func Producer() contract.Producer {
	return queueProducer
}
