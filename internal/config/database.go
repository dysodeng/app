package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type database struct {
	Main db `mapstructure:"main"`
}

type db struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxConnLifetime int    `mapstructure:"max_conn_lifetime"`
}

var Database *database

func databaseConfigLoad() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("database")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	d := v.Sub("database")
	_ = d.BindEnv("main.host", "MAIN_DB_HOST")
	_ = d.BindEnv("main.port", "MAIN_DB_PORT")
	_ = d.BindEnv("main.database", "MAIN_DB_DATABASE")
	_ = d.BindEnv("main.username", "MAIN_DB_USERNAME")
	_ = d.BindEnv("main.password", "MAIN_DB_PASSWORD")
	d.SetDefault("main.host", "127.0.0.1")
	d.SetDefault("main.port", "3306")
	if err := d.Unmarshal(&Database); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/app.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/database.yaml`已变更")
		d = v.Sub("database")
		_ = d.Unmarshal(&Database)
	})
}
