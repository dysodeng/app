package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type redis struct {
	Main  redisItem `mapstructure:"main"`
	Cache redisItem `mapstructure:"cache"`
	MQ    redisItem `mapstructure:"mq"`
}

type redisItem struct {
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	DB        int    `mapstructure:"db"`
	KeyPrefix string `mapstructure:"key_prefix"`
}

var Redis *redis

func redisConfigLoad() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("redis")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	d := v.Sub("redis")
	_ = d.BindEnv("main.host", "MAIN_REDIS_HOST")
	_ = d.BindEnv("main.port", "MAIN_REDIS_PORT")
	_ = d.BindEnv("main.password", "MAIN_REDIS_PASSWORD")
	_ = d.BindEnv("main.database", "MAIN_REDIS_DB")
	_ = d.BindEnv("main.key_prefix", "MAIN_REDIS_KEY_PREFIX")
	_ = d.BindEnv("cache.host", "CACHE_REDIS_HOST")
	_ = d.BindEnv("cache.port", "CACHE_REDIS_PORT")
	_ = d.BindEnv("cache.password", "CACHE_REDIS_PASSWORD")
	_ = d.BindEnv("cache.database", "CACHE_REDIS_DB")
	_ = d.BindEnv("cache.key_prefix", "CACHE_REDIS_KEY_PREFIX")
	_ = d.BindEnv("mq.host", "MQ_REDIS_HOST")
	_ = d.BindEnv("mq.port", "MQ_REDIS_PORT")
	_ = d.BindEnv("mq.password", "MQ_REDIS_PASSWORD")
	_ = d.BindEnv("mq.database", "MQ_REDIS_DB")
	_ = d.BindEnv("mq.key_prefix", "MQ_REDIS_KEY_PREFIX")
	d.SetDefault("main.host", "127.0.0.1")
	d.SetDefault("main.port", "6379")
	d.SetDefault("main.database", 0)
	d.SetDefault("cache.host", "127.0.0.1")
	d.SetDefault("cache.port", "6379")
	d.SetDefault("cache.database", 0)
	d.SetDefault("mq.host", "127.0.0.1")
	d.SetDefault("mq.port", "6379")
	d.SetDefault("mq.database", 0)
	if err := d.Unmarshal(&Redis); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/redis.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/redis.yaml`已变更")
		d = v.Sub("redis")
		_ = d.Unmarshal(&Redis)
	})
}
