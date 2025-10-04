package config

import "github.com/spf13/viper"

type Storage struct {
	Driver    string       `mapstructure:"driver"`
	CdnDomain string       `mapstructure:"cdn_domain"`
	Local     local        `mapstructure:"local"`
	MinIO     cloudStorage `mapstructure:"minio"`
	AliOss    cloudStorage `mapstructure:"ali_oss"`
	HwObs     cloudStorage `mapstructure:"hw_obs"`
	TxCos     cloudStorage `mapstructure:"tx_cos"`
	S3        cloudStorage `mapstructure:"s3"`
}

type local struct {
	RootPath         string `mapstructure:"root_path"`
	MultipartStorage string `mapstructure:"multipart_storage"`
	StaticEnabled    bool   `mapstructure:"static_enabled"`
}

type cloudStorage struct {
	AccessKey            string `mapstructure:"access_key"`
	AccessSecret         string `mapstructure:"access_secret"`
	Endpoint             string `mapstructure:"endpoint"`
	InternalEndpoint     string `mapstructure:"internal_endpoint"`
	WithInternalEndpoint bool   `mapstructure:"with_internal_endpoint"`
	Bucket               string `mapstructure:"bucket"`
	Region               string `mapstructure:"region"`
	UseSSL               bool   `mapstructure:"use_ssl"`
	AccessMode           string `mapstructure:"access_mode"`
}

func storageBindEnv(d *viper.Viper) {
	_ = d.BindEnv("driver", "STORAGE_DRIVER")
	_ = d.BindEnv("cdn_domain", "STORAGE_CDN_DOMAIN")
	_ = d.BindEnv("local.root_path", "LOCAL_ROOT_PATH")
	_ = d.BindEnv("local.multipart_storage", "LOCAL_MULTIPART_STORAGE")
	_ = d.BindEnv("local.static_enabled", "LOCAL_STATIC_ENABLED")
	_ = d.BindEnv("minio.access_key", "MINIO_ACCESS_KEY")
	_ = d.BindEnv("minio.access_secret", "MINIO_SECRET_KEY")
	_ = d.BindEnv("minio.bucket", "MINIO_BUCKET")
	_ = d.BindEnv("minio.endpoint", "MINIO_ENDPOINT")
	_ = d.BindEnv("minio.region", "MINIO_REGION")
	_ = d.BindEnv("minio.use_ssl", "MINIO_USE_SSL")
	_ = d.BindEnv("minio.access_mode", "MINIO_ACCESS_MODE")
	_ = d.BindEnv("ali_oss.access_key", "ALI_OSS_ACCESS_ID")
	_ = d.BindEnv("ali_oss.access_secret", "ALI_OSS_ACCESS_KEY")
	_ = d.BindEnv("ali_oss.bucket", "ALI_OSS_BUCKET")
	_ = d.BindEnv("ali_oss.endpoint", "ALI_OSS_ENDPOINT")
	_ = d.BindEnv("ali_oss.endpoint_internal", "ALI_OSS_INTERNAL_ENDPOINT")
	_ = d.BindEnv("ali_oss.with_internal_endpoint", "ALI_OSS_WITH_INTERNAL_ENDPOINT")
	_ = d.BindEnv("ali_oss.region", "ALI_OSS_REGION")
	_ = d.BindEnv("ali_oss.access_mode", "ALI_OSS_ACCESS_MODE")
	_ = d.BindEnv("hw_obs.access_key", "HW_OBS_ACCESS_KEY")
	_ = d.BindEnv("hw_obs.access_secret", "HW_OBS_SECRET_KEY")
	_ = d.BindEnv("hw_obs.bucket", "HW_OBS_BUCKET")
	_ = d.BindEnv("hw_obs.endpoint", "HW_OBS_BUCKET")
	_ = d.BindEnv("hw_obs.endpoint_internal", "HW_OBS_INTERNAL_ENDPOINT")
	_ = d.BindEnv("hw_obs.with_internal_endpoint", "HW_OBS_WITH_INTERNAL_ENDPOINT")
	_ = d.BindEnv("hw_obs.access_mode", "HW_OBS_ACCESS_MODE")
	_ = d.BindEnv("tx_cos.access_key", "TX_COS_ACCESS_KEY")
	_ = d.BindEnv("tx_cos.access_secret", "TX_COS_SECRET_KEY")
	_ = d.BindEnv("tx_cos.bucket", "TX_COS_BUCKET")
	_ = d.BindEnv("tx_cos.endpoint", "TX_COS_ENDPOINT")
	_ = d.BindEnv("tx_cos.endpoint_internal", "TX_COS_INTERNAL_ENDPOINT")
	_ = d.BindEnv("tx_cos.with_internal_endpoint", "TX_COS_WITH_INTERNAL_ENDPOINT")
	_ = d.BindEnv("tx_cos.access_mode", "TX_COS_ACCESS_MODE")
	_ = d.BindEnv("s3.access_key", "S3_ACCESS_KEY")
	_ = d.BindEnv("s3.access_secret", "S3_SECRET_KEY")
	_ = d.BindEnv("s3.bucket", "S3_BUCKET")
	_ = d.BindEnv("s3.endpoint", "S3_ENDPOINT")
	_ = d.BindEnv("s3.endpoint_internal", "S3_INTERNAL_ENDPOINT")
	_ = d.BindEnv("s3.with_internal_endpoint", "S3_WITH_INTERNAL_ENDPOINT")
	_ = d.BindEnv("s3.access_mode", "S3_ACCESS_MODE")
	d.SetDefault("driver", "local")
	d.SetDefault("minio.access_mode", "private")
	d.SetDefault("ali_oss.access_mode", "private")
	d.SetDefault("hw_obs.access_mode", "private")
	d.SetDefault("tx_cos.access_mode", "private")
	d.SetDefault("s3.access_mode", "private")
}
