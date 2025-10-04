package config

import (
	"time"

	"github.com/spf13/viper"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Migration       Migration     `mapstructure:"migration"`
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type Migration struct {
	Enabled bool `mapstructure:"enabled"`
}

func databaseBindEnv(v *viper.Viper) {
	_ = v.BindEnv("host", "MAIN_DB_HOST")
	_ = v.BindEnv("port", "MAIN_DB_PORT")
	_ = v.BindEnv("database", "MAIN_DB_DATABASE")
	_ = v.BindEnv("username", "MAIN_DB_USERNAME")
	_ = v.BindEnv("password", "MAIN_DB_PASSWORD")
	_ = v.BindEnv("tracer_enable", "MAIN_DB_TRACER_ENABLE")
	v.SetDefault("host", "127.0.0.1")
	v.SetDefault("port", "3306")
}
