package config

import (
	"github.com/spf13/viper"
)

// Cache 应用缓存
type Cache struct {
	Driver string `mapstructure:"driver"`
}

func cacheBindEnv(d *viper.Viper) {
	d.SetDefault("driver", "memory")
}
