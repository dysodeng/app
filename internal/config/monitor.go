package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Monitor 平台监控配置
type monitor struct {
	Tracer tracer `mapstructure:"tracer"`
}

// tracer 链路追踪配置
type tracer struct {
	ReportEnabled  bool   `mapstructure:"report_enabled"`
	ReportEndpoint string `mapstructure:"report_endpoint"`
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
}

var Monitor *monitor

func monitorConfigLoad() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("monitor")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	mon := v.Sub("monitor")
	_ = mon.BindEnv("tracer.report_enabled", "MONITOR_TRACER_REPORT_ENABLED")
	_ = mon.BindEnv("tracer.report_endpoint", "MONITOR_TRACER_REPORT_ENDPOINT")
	_ = mon.BindEnv("tracer.service_name", "MONITOR_TRACER_SERVICE_NAME")
	_ = mon.BindEnv("tracer.service_version", "MONITOR_TRACER_SERVICE_VERSION")
	mon.SetDefault("tracer.report_enabled", "false")
	mon.SetDefault("tracer.report_endpoint", "http://127.0.0.1:4318")
	mon.SetDefault("tracer.service_name", "app")
	mon.SetDefault("tracer.service_version", "v1.0.0")

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
