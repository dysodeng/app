package config

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
