package storage

import "github.com/dysodeng/fs/driver/local"

// AccessMode 访问模式
type AccessMode string

const (
	Private         AccessMode = "private"           // 私有读写
	PublicRead      AccessMode = "public-read"       // 公共读
	PublicReadWrite AccessMode = "public-read-write" // 公共读写
)

type Config struct {
	Driver    string
	CdnDomain string       `json:"cdn_domain"`
	Local     Local        `json:"local"`
	Minio     CloudStorage `json:"minio"`
	AliOss    CloudStorage `json:"ali_oss"`
	HwObs     CloudStorage `json:"hw_obs"`
	TxCos     CloudStorage `json:"tx_cos"`
	S3        CloudStorage `json:"s3"`
}

type CloudStorage struct {
	AccessKey            string     `json:"access_key"`
	AccessSecret         string     `json:"access_secret"`
	Bucket               string     `json:"bucket"`
	Endpoint             string     `json:"endpoint"`
	InternalEndpoint     string     `json:"internal_endpoint"`
	WithInternalEndpoint bool       `json:"with_internal_endpoint"`
	Region               string     `json:"region"`
	UseSSL               bool       `json:"use_ssl"`
	AccessMode           AccessMode `json:"access_mode"`
}

type Local struct {
	RootPath         string                 `json:"root_path"`
	MultipartStorage local.MultipartStorage `json:"multipart_storage"`
}
