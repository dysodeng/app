package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// LoadResult 配置加载结果
type LoadResult struct {
	Config      *Config
	Source      string // "local" 或配置源类型如 "etcd"
	SourcePath  string
	LocalPath   string
	WatchSource Source // 仅配置中心模式下非 nil
}

// Load 按环境变量、配置中心、本地文件的优先级加载配置。
func Load(path string) (*LoadResult, error) {
	return load(path, createSource)
}

func load(path string, sourceFactory func(*ConfigCenterConfig) Source) (*LoadResult, error) {
	_ = godotenv.Load()

	v, err := readLocal(path)
	if err != nil {
		return nil, err
	}
	bootstrapCenter, err := decodeConfigCenter(v)
	if err != nil {
		return nil, fmt.Errorf("解析引导配置失败: %w", err)
	}

	result := &LoadResult{Source: "local", SourcePath: path, LocalPath: path}
	if source := sourceFactory(bootstrapCenter); source != nil {
		result.WatchSource = source
		data, err := source.Load()
		if err != nil {
			slog.Warn("从配置中心加载失败，回退使用本地配置",
				"source", source.Type(), "error", err)
		} else if err := v.MergeConfig(bytes.NewReader(data)); err != nil {
			slog.Warn("配置中心内容解析失败，回退使用本地配置",
				"source", source.Type(), "error", err)
		} else {
			result.Source = source.Type()
			result.SourcePath = sourceKey(bootstrapCenter)
		}
	}

	cfg, err := decodeConfig(v)
	if err != nil {
		return nil, err
	}
	Publish(cfg)
	result.Config = cfg
	return result, nil
}

func readLocal(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	return v, nil
}

// decodeConfig 在合并本地和配置中心后应用环境变量。
func decodeConfig(v *viper.Viper) (*Config, error) {
	var appConfig AppConfig
	app := section(v, "app")
	appBindEnv(app)
	if err := app.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	var serverConfig Server
	server := section(v, "server")
	serverBindEnv(server)
	if err := server.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}

	var securityConfig Security
	security := section(v, "security")
	securityBindEnv(security)
	if err := security.Unmarshal(&securityConfig); err != nil {
		return nil, err
	}

	var databaseConfig DatabaseConfig
	database := section(v, "database")
	databaseBindEnv(database)
	if err := database.Unmarshal(&databaseConfig); err != nil {
		return nil, err
	}

	var redisConfig Redis
	redis := section(v, "redis")
	redisBindEnv(redis)
	if err := redis.Unmarshal(&redisConfig); err != nil {
		return nil, err
	}

	var cacheConfig Cache
	cache := section(v, "cache")
	cacheBindEnv(cache)
	if err := cache.Unmarshal(&cacheConfig); err != nil {
		return nil, err
	}

	var messageQueueConfig MessageQueue
	messageQueue := section(v, "message_queue")
	messageQueueBindEnv(messageQueue)
	if err := messageQueue.Unmarshal(&messageQueueConfig); err != nil {
		return nil, err
	}

	var etcdConfig Etcd
	etcd := section(v, "etcd")
	etcdBindEnv(etcd)
	if err := etcd.Unmarshal(&etcdConfig); err != nil {
		return nil, err
	}

	var storageConfig Storage
	storage := section(v, "storage")
	storageBindEnv(storage)
	if err := storage.Unmarshal(&storageConfig); err != nil {
		return nil, err
	}

	var monitorConfig Monitor
	monitor := section(v, "monitor")
	monitorBindEnv(monitor)
	if err := monitor.Unmarshal(&monitorConfig); err != nil {
		return nil, err
	}

	var thirdPartyConfig ThirdParty
	thirdParty := section(v, "third_party")
	thirdPartyBindEnv(thirdParty)
	if err := thirdParty.Unmarshal(&thirdPartyConfig); err != nil {
		return nil, err
	}

	configCenter, err := decodeConfigCenter(v)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		ConfigCenter: configCenter,
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

	return cfg, nil
}

func section(v *viper.Viper, name string) *viper.Viper {
	sub := v.Sub(name)
	if sub == nil {
		sub = viper.New()
	}
	sub.AllowEmptyEnv(true)
	return sub
}

func decodeConfigCenter(v *viper.Viper) (*ConfigCenterConfig, error) {
	center := section(v, "config_center")
	_ = center.BindEnv("type", "CONFIG_CENTER_TYPE")
	_ = center.BindEnv("etcd.endpoints", "CONFIG_CENTER_ETCD_ENDPOINTS")
	_ = center.BindEnv("etcd.key", "CONFIG_CENTER_ETCD_KEY")
	_ = center.BindEnv("etcd.timeout", "CONFIG_CENTER_ETCD_TIMEOUT")
	_ = center.BindEnv("etcd.username", "CONFIG_CENTER_ETCD_USERNAME")
	_ = center.BindEnv("etcd.password", "CONFIG_CENTER_ETCD_PASSWORD")
	if value, ok := os.LookupEnv("CONFIG_CENTER_ETCD_ENDPOINTS"); ok {
		if value == "" {
			center.Set("etcd.endpoints", []string{})
		} else {
			center.Set("etcd.endpoints", strings.Split(value, ","))
		}
	}
	var cfg ConfigCenterConfig
	if err := center.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if cfg.Type == "" && cfg.Etcd == nil {
		return nil, nil
	}
	return &cfg, nil
}

// createSource 根据配置中心配置创建配置源
func createSource(cc *ConfigCenterConfig) Source {
	if cc == nil || cc.Type == "" || cc.Type == "static" {
		return nil
	}
	switch cc.Type {
	case "etcd":
		if cc.Etcd == nil {
			return nil
		}
		return NewEtcdSource(cc.Etcd)
	default:
		slog.Warn("不支持的配置中心类型", "type", cc.Type)
		return nil
	}
}

// sourceKey 获取配置源的 key 信息
func sourceKey(cc *ConfigCenterConfig) string {
	if cc == nil {
		return ""
	}
	if cc.Type == "etcd" && cc.Etcd != nil {
		return cc.Etcd.Key
	}
	return ""
}
