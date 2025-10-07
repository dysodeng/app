package config

import "github.com/spf13/viper"

// Server 服务配置
type Server struct {
	HTTP      HTTPConfig      `mapstructure:"http"`
	GRPC      GRPCConfig      `mapstructure:"grpc"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
	Event     EventConfig     `mapstructure:"event"`
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

// GRPCConfig gRPC配置
type GRPCConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
}

// WebSocketConfig WebSocket配置
type WebSocketConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

// EventConfig 事件消费者服务配置
type EventConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Driver  string `mapstructure:"driver"`
}

func serverBindEnv(v *viper.Viper) {
	_ = v.BindEnv("http.port", "SERVER_HTTP_PORT")
	_ = v.BindEnv("grpc.port", "SERVER_GRPC_PORT")
	_ = v.BindEnv("websocket.port", "SERVER_WEBSOCKET_PORT")
}
