package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type mq struct {
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

var MQ *mq

func mqConfigLoad() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("mq")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	d := v.Sub("mq")
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

	if err := d.Unmarshal(&MQ); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/mq.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/mq.yaml`已变更")
		d = v.Sub("mq")
		_ = d.Unmarshal(&MQ)
	})
}
