package config

import (
	"github.com/spf13/viper"
)

// Cache 应用缓存
type Cache struct {
	Driver     string `mapstructure:"driver"`
	Serializer string `mapstructure:"serializer"`
}

func cacheBindEnv(d *viper.Viper) {
	d.SetDefault("driver", "memory")
	d.SetDefault("serializer", "json")
}
