package config

import "github.com/spf13/viper"

// Server 服务配置
type Server struct {
	HTTP      HTTPConfig      `mapstructure:"http"`
	GRPC      GRPCConfig      `mapstructure:"grpc"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

// GRPCConfig gRPC配置
type GRPCConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

// WebSocketConfig WebSocket配置
type WebSocketConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

func serverBindEnv(v *viper.Viper) {
	_ = v.BindEnv("http.port", "SERVER_HTTP_PORT")
	_ = v.BindEnv("grpc.port", "SERVER_GRPC_PORT")
	_ = v.BindEnv("websocket.port", "SERVER_WEBSOCKET_PORT")
}
