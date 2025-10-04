package config

import "github.com/spf13/viper"

// Redis redis配置
type Redis struct {
	Main  RedisItem `mapstructure:"main"`
	Cache RedisItem `mapstructure:"cache"`
	MQ    RedisItem `mapstructure:"mq"`
}

type RedisItem struct {
	// 单机模式配置
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	DB        int    `mapstructure:"db"`
	KeyPrefix string `mapstructure:"key_prefix"`

	// 模式选择: standalone, cluster, sentinel
	Mode string `mapstructure:"mode"`

	// 集群模式配置
	Cluster ClusterConfig `mapstructure:"cluster"`

	// 哨兵模式配置
	Sentinel SentinelConfig `mapstructure:"sentinel"`

	// 连接池配置
	Pool PoolConfig `mapstructure:"pool"`
}

type ClusterConfig struct {
	Addrs    []string `mapstructure:"addrs"`
	Password string   `mapstructure:"password"`
}

type SentinelConfig struct {
	MasterName       string   `mapstructure:"master_name"`
	SentinelAddrs    []string `mapstructure:"sentinel_addrs"`
	Password         string   `mapstructure:"password"`
	SentinelPassword string   `mapstructure:"sentinel_password"`
	DB               int      `mapstructure:"db"`
}

type PoolConfig struct {
	MinIdleConns int `mapstructure:"min_idle_conns"`
	MaxRetries   int `mapstructure:"max_retries"`
	PoolSize     int `mapstructure:"pool_size"`
}

func redisBindEnv(d *viper.Viper) {
	// 单机模式环境变量绑定
	_ = d.BindEnv("main.host", "MAIN_REDIS_HOST")
	_ = d.BindEnv("main.port", "MAIN_REDIS_PORT")
	_ = d.BindEnv("main.password", "MAIN_REDIS_PASSWORD")
	_ = d.BindEnv("main.database", "MAIN_REDIS_DB")
	_ = d.BindEnv("main.key_prefix", "MAIN_REDIS_KEY_PREFIX")
	_ = d.BindEnv("main.mode", "MAIN_REDIS_MODE")

	// 集群模式环境变量绑定
	_ = d.BindEnv("main.cluster.addrs", "MAIN_REDIS_CLUSTER_ADDRS")
	_ = d.BindEnv("main.cluster.password", "MAIN_REDIS_CLUSTER_PASSWORD")

	// 哨兵模式环境变量绑定
	_ = d.BindEnv("main.sentinel.master_name", "MAIN_REDIS_SENTINEL_MASTER_NAME")
	_ = d.BindEnv("main.sentinel.sentinel_addrs", "MAIN_REDIS_SENTINEL_ADDRS")
	_ = d.BindEnv("main.sentinel.password", "MAIN_REDIS_SENTINEL_PASSWORD")
	_ = d.BindEnv("main.sentinel.sentinel_password", "MAIN_REDIS_SENTINEL_SENTINEL_PASSWORD")
	_ = d.BindEnv("main.sentinel.db", "MAIN_REDIS_SENTINEL_DB")

	// Cache和MQ的类似绑定...
	_ = d.BindEnv("cache.host", "CACHE_REDIS_HOST")
	_ = d.BindEnv("cache.port", "CACHE_REDIS_PORT")
	_ = d.BindEnv("cache.password", "CACHE_REDIS_PASSWORD")
	_ = d.BindEnv("cache.database", "CACHE_REDIS_DB")
	_ = d.BindEnv("cache.key_prefix", "CACHE_REDIS_KEY_PREFIX")
	_ = d.BindEnv("cache.mode", "CACHE_REDIS_MODE")
	_ = d.BindEnv("mq.host", "MQ_REDIS_HOST")
	_ = d.BindEnv("mq.port", "MQ_REDIS_PORT")
	_ = d.BindEnv("mq.password", "MQ_REDIS_PASSWORD")
	_ = d.BindEnv("mq.database", "MQ_REDIS_DB")
	_ = d.BindEnv("mq.key_prefix", "MQ_REDIS_KEY_PREFIX")
	_ = d.BindEnv("mq.mode", "MQ_REDIS_MODE")

	// 默认值设置
	d.SetDefault("main.host", "127.0.0.1")
	d.SetDefault("main.port", "6379")
	d.SetDefault("main.database", 0)
	d.SetDefault("main.mode", "standalone")
	d.SetDefault("main.pool.min_idle_conns", 10)
	d.SetDefault("main.pool.max_retries", 3)
	d.SetDefault("main.pool.pool_size", 100)

	d.SetDefault("cache.host", "127.0.0.1")
	d.SetDefault("cache.port", "6379")
	d.SetDefault("cache.database", 0)
	d.SetDefault("cache.mode", "standalone")
	d.SetDefault("cache.pool.min_idle_conns", 10)
	d.SetDefault("cache.pool.max_retries", 3)
	d.SetDefault("cache.pool.pool_size", 100)

	d.SetDefault("mq.host", "127.0.0.1")
	d.SetDefault("mq.port", "6379")
	d.SetDefault("mq.database", 0)
	d.SetDefault("mq.mode", "standalone")
	d.SetDefault("mq.pool.min_idle_conns", 10)
	d.SetDefault("mq.pool.max_retries", 3)
	d.SetDefault("mq.pool.pool_size", 100)
}
