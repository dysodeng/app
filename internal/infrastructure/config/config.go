package config

import (
	"github.com/joho/godotenv"
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
	App          AppConfig      `mapstructure:"app"`
	Server       Server         `mapstructure:"server"`
	Security     Security       `mapstructure:"security"`
	Database     DatabaseConfig `mapstructure:"database"`
	Redis        Redis          `mapstructure:"redis"`
	Cache        Cache          `mapstructure:"cache"`
	MessageQueue MessageQueue   `mapstructure:"message_queue"`
	Etcd         Etcd           `mapstructure:"etcd"`
	Storage      Storage        `mapstructure:"storage"`
	Monitor      Monitor        `mapstructure:"monitor"`
	ThirdParty   ThirdParty     `mapstructure:"third_party"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 加载.env
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigFile(configPath)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var appConfig AppConfig
	app := v.Sub("app")
	appBindEnv(app)
	if err := app.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	var serverConfig Server
	server := v.Sub("server")
	serverBindEnv(server)
	if err := server.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}

	var securityConfig Security
	security := v.Sub("security")
	securityBindEnv(security)
	if err := security.Unmarshal(&securityConfig); err != nil {
		return nil, err
	}

	var databaseConfig DatabaseConfig
	database := v.Sub("database")
	databaseBindEnv(database)
	if err := database.Unmarshal(&databaseConfig); err != nil {
		return nil, err
	}

	var redisConfig Redis
	redis := v.Sub("redis")
	redisBindEnv(redis)
	if err := redis.Unmarshal(&redisConfig); err != nil {
		return nil, err
	}

	var cacheConfig Cache
	cache := v.Sub("cache")
	cacheBindEnv(cache)
	if err := cache.Unmarshal(&cacheConfig); err != nil {
		return nil, err
	}

	var messageQueueConfig MessageQueue
	messageQueue := v.Sub("message_queue")
	messageQueueBindEnv(messageQueue)
	if err := messageQueue.Unmarshal(&messageQueueConfig); err != nil {
		return nil, err
	}

	var etcdConfig Etcd
	etcd := v.Sub("etcd")
	etcdBindEnv(etcd)
	if err := etcd.Unmarshal(&etcdConfig); err != nil {
		return nil, err
	}

	var storageConfig Storage
	storage := v.Sub("storage")
	storageBindEnv(storage)
	if err := storage.Unmarshal(&storageConfig); err != nil {
		return nil, err
	}

	var monitorConfig Monitor
	monitor := v.Sub("monitor")
	monitorBindEnv(monitor)
	if err := monitor.Unmarshal(&monitorConfig); err != nil {
		return nil, err
	}

	var thirdPartyConfig ThirdParty
	thirdParty := v.Sub("third_party")
	thirdPartyBindEnv(thirdParty)
	if err := thirdParty.Unmarshal(&thirdPartyConfig); err != nil {
		return nil, err
	}

	config := Config{
		App:          appConfig,
		Server:       serverConfig,
		Security:     securityConfig,
		Database:     databaseConfig,
		Redis:        redisConfig,
		Cache:        cacheConfig,
		MessageQueue: messageQueueConfig,
		Etcd:         etcdConfig,
		Storage:      storageConfig,
		Monitor:      monitorConfig,
		ThirdParty:   thirdPartyConfig,
	}

	GlobalConfig = &config

	return &config, nil
}
