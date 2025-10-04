package config

import (
	"github.com/spf13/viper"
)

const (
	VarPath  string = "var"
	LogPath         = VarPath + "/logs"
	TempPath        = VarPath + "/tmp"
)

var GlobalConfig *Config

// Config 应用配置
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   Server         `mapstructure:"server"`
	Security Security       `mapstructure:"security"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    Redis          `mapstructure:"redis"`
}

// AppConfig 应用基本配置
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// Security 安全配置
type Security struct {
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	GlobalConfig = &config

	return &config, nil
}
