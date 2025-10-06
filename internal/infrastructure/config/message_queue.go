package config

import (
	"github.com/spf13/viper"
)

// MessageQueue 消息队列配置
type MessageQueue struct {
	Enabled bool   `mapstructure:"enabled"`
	Driver  string `mapstructure:"driver"`
	Amqp    amqp   `mapstructure:"amqp"`
}

type amqp struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Vhost    string `mapstructure:"vhost"`
}

func messageQueueBindEnv(d *viper.Viper) {
	_ = d.BindEnv("driver", "MQ_DRIVER")
	_ = d.BindEnv("amqp.host", "MQ_AMQP_HOST")
	_ = d.BindEnv("amqp.port", "MQ_AMQP_PORT")
	_ = d.BindEnv("amqp.username", "MQ_AMQP_USERNAME")
	_ = d.BindEnv("amqp.password", "MQ_AMQP_PASSWORD")
	_ = d.BindEnv("amqp.password", "MQ_AMQP_PASSWORD")
	_ = d.BindEnv("amqp.vhost", "MQ_AMQP_VHOST")
	_ = d.BindEnv("redis.connection", "MQ_REDIS_CONNECTION")
	d.SetDefault("driver", "amqp")
	d.SetDefault("amqp.host", "127.0.0.1")
	d.SetDefault("amqp.port", "5672")
	d.SetDefault("amqp.username", "guest")
	d.SetDefault("amqp.password", "guest")
	d.SetDefault("amqp.vhost", "/")
}
