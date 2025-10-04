package config

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
