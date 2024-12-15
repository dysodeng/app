package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type filesystem struct {
	Storage string `mapstructure:"storage"`
	Local   struct {
		BasePath  string `mapstructure:"base_path"`
		LogicPath string `mapstructure:"logic_path"`
		BaseUrl   string `mapstructure:"base_url"`
	} `mapstructure:"local"`
	AliOss cloudStorage `mapstructure:"ali_oss"`
	HwObs  cloudStorage `mapstructure:"hw_obs"`
	Minio  cloudStorage `mapstructure:"minio"`
}

type cloudStorage struct {
	AccessKey        string `mapstructure:"access_key"`
	AccessSecret     string `mapstructure:"access_secret"`
	Bucket           string `mapstructure:"bucket"`
	Endpoint         string `mapstructure:"endpoint"`
	EndpointInternal string `mapstructure:"endpoint_internal"`
	Region           string `mapstructure:"region"`
	StsRoleArn       string `mapstructure:"sts_role_arn"`
	CdnDomain        string `mapstructure:"cdn_domain"`
	UseSSL           bool   `mapstructure:"use_ssl"`
}

var Filesystem *filesystem

func loadFilesystemConfig() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("filesystem")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	d := v.Sub("filesystem")
	_ = d.BindEnv("storage", "FILESYSTEM_STORAGE")
	_ = d.BindEnv("local.base_path", "LOCAL_BASE_PATH")
	_ = d.BindEnv("local.logic_path", "LOCAL_LOGIC_PATH")
	_ = d.BindEnv("local.base_url", "LOCAL_BASE_URL")
	_ = d.BindEnv("ali_oss.access_key", "ALI_OSS_ACCESS_ID")
	_ = d.BindEnv("ali_oss.access_secret", "ALI_OSS_ACCESS_KEY")
	_ = d.BindEnv("ali_oss.bucket", "ALI_OSS_BUCKET")
	_ = d.BindEnv("ali_oss.endpoint", "ALI_OSS_ENDPOINT")
	_ = d.BindEnv("ali_oss.endpoint_internal", "ALI_OSS_ENDPOINT_INTERNAL")
	_ = d.BindEnv("ali_oss.region", "ALI_OSS_REGION")
	_ = d.BindEnv("ali_oss.sts_role_arn", "ALI_OSS_STS_ROLE_ARN")
	_ = d.BindEnv("ali_oss.cdn_domain", "ALI_OSS_CDN_DOMAIN")
	_ = d.BindEnv("hw_obs.access_key", "HW_OBS_ACCESS_KEY")
	_ = d.BindEnv("hw_obs.access_secret", "HW_OBS_SECRET_KEY")
	_ = d.BindEnv("hw_obs.bucket", "HW_OBS_BUCKET")
	_ = d.BindEnv("hw_obs.endpoint", "HW_OBS_ENDPOINT")
	_ = d.BindEnv("hw_obs.cdn_domain", "HW_OBS_CDN_DOMAIN")
	_ = d.BindEnv("minio.access_key", "MINIO_ACCESS_KEY")
	_ = d.BindEnv("minio.access_secret", "MINIO_SECRET_KEY")
	_ = d.BindEnv("minio.bucket", "MINIO_BUCKET")
	_ = d.BindEnv("minio.endpoint", "MINIO_ENDPOINT")
	_ = d.BindEnv("minio.cdn_domain", "MINIO_CDN_DOMAIN")
	d.SetDefault("storage", "local")
	if err := d.Unmarshal(&Filesystem); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/filesystem.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/filesystem.yaml`已变更")
		d = v.Sub("filesystem")
		_ = d.Unmarshal(&Filesystem)
	})
}
