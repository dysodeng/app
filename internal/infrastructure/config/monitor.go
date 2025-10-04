package config

// Monitor 平台监控配置
type Monitor struct {
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
	Tracer         otlp   `mapstructure:"tracer"`
	Metrics        otlp   `mapstructure:"metrics"`
	Log            otlp   `mapstructure:"log"`
}

// tracer 链路追踪配置
type otlp struct {
	OtlpEnabled  bool   `mapstructure:"otlp_enabled"`
	OtlpEndpoint string `mapstructure:"otlp_endpoint"`
}
