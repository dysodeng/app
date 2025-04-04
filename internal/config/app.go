package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type appConfig struct {
	Env    Env    `mapstructure:"env"`
	Debug  bool   `mapstructure:"debug"`
	Name   string `mapstructure:"name"`
	Domain string `mapstructure:"domain"`
	Jwt    struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

type serverConfig struct {
	Http struct {
		Enabled bool   `mapstructure:"enabled"`
		Port    string `mapstructure:"port"`
	} `mapstructure:"http"`
	Grpc struct {
		Enabled bool   `mapstructure:"enabled"`
		Port    string `mapstructure:"port"`
	} `mapstructure:"grpc"`
	Task struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"task"`
	Cron struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"cron"`
	Websocket struct {
		Enabled bool   `mapstructure:"enabled"`
		Port    string `mapstructure:"port"`
	} `mapstructure:"websocket"`
	Health struct {
		Enabled bool   `mapstructure:"enabled"`
		Port    string `mapstructure:"port"`
	} `mapstructure:"health"`
}

var App *appConfig
var Server *serverConfig

func loadAppConfig() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	app := v.Sub("app")
	_ = app.BindEnv("env", "APP_ENV")
	_ = app.BindEnv("env", "APP_DEBUG")
	_ = app.BindEnv("domain", "APP_DOMAIN")
	_ = app.BindEnv("jwt.secret", "APP_JWT_SECRET")
	app.SetDefault("env", "dev")
	if err := app.Unmarshal(&App); err != nil {
		panic(err)
	}

	server := v.Sub("server")
	_ = server.BindEnv("http.port", "SERVER_HTTP_PORT")
	_ = server.BindEnv("grpc.port", "SERVER_GRPC_PORT")
	_ = server.BindEnv("health.port", "SERVER_HEALTH_CHECK_PORT")
	if err := server.Unmarshal(&Server); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/app.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/app.yaml`已变更")
		app = v.Sub("app")
		server = v.Sub("server")
		_ = app.Unmarshal(&App)
		_ = server.Unmarshal(&Server)
	})
}
