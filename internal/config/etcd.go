package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type etcd struct {
	Grpc etcdItem `mapstructure:"grpc"`
}

type etcdItem struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	TLS      bool   `mapstructure:"tls"`
}

var Etcd *etcd

func etcdConfigLoad() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("etcd")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	d := v.Sub("etcd")
	_ = d.BindEnv("grpc.addr", "GRPC_ETCD_ADDR")
	_ = d.BindEnv("grpc.username", "GRPC_ETCD_USERNAME")
	_ = d.BindEnv("grpc.password", "GRPC_ETCD_PASSWORD")
	_ = d.BindEnv("grpc.tls", "GRPC_ETCD_TLS")
	d.SetDefault("grpc.addr", "127.0.0.1:2379")
	if err := d.Unmarshal(&Etcd); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/redis.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/etcd.yaml`已变更")
		d = v.Sub("etcd")
		_ = d.Unmarshal(&Etcd)
	})
}
