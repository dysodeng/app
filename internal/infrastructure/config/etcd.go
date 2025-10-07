package config

import (
	"github.com/spf13/viper"
)

type Etcd struct {
	GRPC etcdItem `mapstructure:"grpc"`
}

type etcdItem struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	TLS      bool   `mapstructure:"tls"`
}

func etcdBindEnv(d *viper.Viper) {
	_ = d.BindEnv("grpc.addr", "GRPC_ETCD_ADDR")
	_ = d.BindEnv("grpc.username", "GRPC_ETCD_USERNAME")
	_ = d.BindEnv("grpc.password", "GRPC_ETCD_PASSWORD")
	_ = d.BindEnv("grpc.tls", "GRPC_ETCD_TLS")
	d.SetDefault("grpc.addr", "127.0.0.1:2379")
}
