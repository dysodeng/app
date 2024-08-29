package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// cache 应用缓存
type cache struct {
	Driver string `mapstructure:"driver"`
}

var Cache *cache

func cacheConfigLoad() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("cache")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	d := v.Sub("cache")
	d.SetDefault("driver", "memory")
	if err := d.Unmarshal(&Cache); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/cache.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/cache.yaml`已变更")
		d = v.Sub("cache")
		_ = d.Unmarshal(&Cache)
	})
}
