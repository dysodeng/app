package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Monitor 平台监控配置
type monitor struct {
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

var Monitor *monitor

func loadMonitorConfig() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("monitor")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	mon := v.Sub("monitor")
	_ = mon.BindEnv("service_name", "MONITOR_SERVICE_NAME")
	_ = mon.BindEnv("service_version", "MONITOR_SERVICE_VERSION")
	_ = mon.BindEnv("tracer.otlp_enabled", "MONITOR_TRACER_OTLP_ENABLED")
	_ = mon.BindEnv("tracer.otlp_endpoint", "MONITOR_TRACER_OTLP_ENDPOINT")
	_ = mon.BindEnv("metrics.otlp_enabled", "MONITOR_METRICS_OTLP_ENABLED")
	_ = mon.BindEnv("metrics.otlp_endpoint", "MONITOR_METRICS_OTLP_ENDPOINT")
	_ = mon.BindEnv("log.otlp_enabled", "MONITOR_LOG_OTLP_ENABLED")
	_ = mon.BindEnv("log.otlp_endpoint", "MONITOR_LOG_OTLP_ENDPOINT")
	mon.SetDefault("service_name", "app")
	mon.SetDefault("service_version", "v1.0.0")
	mon.SetDefault("tracer.otlp_enabled", "false")
	mon.SetDefault("tracer.otlp_endpoint", "http://127.0.0.1:4318")
	mon.SetDefault("metrics.otlp_enabled", "false")
	mon.SetDefault("metrics.otlp_endpoint", "http://127.0.0.1:4318")
	mon.SetDefault("log.otlp_enabled", "false")
	mon.SetDefault("log.otlp_endpoint", "127.0.0.1:4318")

	if err := mon.Unmarshal(&Monitor); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/monitor.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/monitor.yaml`已变更")
		mon = v.Sub("monitor")
		_ = mon.Unmarshal(&Monitor)
	})
}
