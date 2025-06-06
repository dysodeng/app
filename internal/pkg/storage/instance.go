package storage

import (
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/redis"
	"github.com/dysodeng/fs/driver/local"
)

var cfg Config

func init() {
	var localMultipartStorage local.MultipartStorage
	if config.Storage.Driver == "local" {
		switch config.Storage.Local.MultipartStorage {
		case "redis":
			localMultipartStorage = NewRedisMultipartStorage(redis.MainClient(), "")
		default:
			// 默认为文件
			localMultipartStorage = nil
		}
	}
	cfg = Config{
		Driver:    config.Storage.Driver,
		CdnDomain: config.Storage.CdnDomain,
		Local: Local{
			RootPath:         config.Storage.Local.RootPath,
			MultipartStorage: localMultipartStorage,
		},
		Minio: CloudStorage{
			AccessKey:            config.Storage.MinIO.AccessKey,
			AccessSecret:         config.Storage.MinIO.AccessSecret,
			Bucket:               config.Storage.MinIO.Bucket,
			Endpoint:             config.Storage.MinIO.Endpoint,
			InternalEndpoint:     config.Storage.MinIO.InternalEndpoint,
			WithInternalEndpoint: config.Storage.MinIO.WithInternalEndpoint,
			Region:               config.Storage.MinIO.Region,
			UseSSL:               config.Storage.MinIO.UseSSL,
			AccessMode:           AccessMode(config.Storage.MinIO.AccessMode),
		},
		HwObs: CloudStorage{
			AccessKey:            config.Storage.HwObs.AccessKey,
			AccessSecret:         config.Storage.HwObs.AccessSecret,
			Bucket:               config.Storage.HwObs.Bucket,
			Endpoint:             config.Storage.HwObs.Endpoint,
			InternalEndpoint:     config.Storage.HwObs.InternalEndpoint,
			WithInternalEndpoint: config.Storage.HwObs.WithInternalEndpoint,
			Region:               config.Storage.HwObs.Region,
			UseSSL:               config.Storage.HwObs.UseSSL,
			AccessMode:           AccessMode(config.Storage.HwObs.AccessMode),
		},
		AliOss: CloudStorage{
			AccessKey:            config.Storage.AliOss.AccessKey,
			AccessSecret:         config.Storage.AliOss.AccessSecret,
			Bucket:               config.Storage.AliOss.Bucket,
			Endpoint:             config.Storage.AliOss.Endpoint,
			InternalEndpoint:     config.Storage.AliOss.InternalEndpoint,
			WithInternalEndpoint: config.Storage.AliOss.WithInternalEndpoint,
			Region:               config.Storage.AliOss.Region,
			UseSSL:               config.Storage.AliOss.UseSSL,
			AccessMode:           AccessMode(config.Storage.AliOss.AccessMode),
		},
		TxCos: CloudStorage{
			AccessKey:            config.Storage.TxCos.AccessKey,
			AccessSecret:         config.Storage.TxCos.AccessSecret,
			Bucket:               config.Storage.TxCos.Bucket,
			Endpoint:             config.Storage.TxCos.Endpoint,
			InternalEndpoint:     config.Storage.TxCos.InternalEndpoint,
			WithInternalEndpoint: config.Storage.TxCos.WithInternalEndpoint,
			Region:               config.Storage.TxCos.Region,
			UseSSL:               config.Storage.TxCos.UseSSL,
			AccessMode:           AccessMode(config.Storage.TxCos.AccessMode),
		},
		S3: CloudStorage{
			AccessKey:            config.Storage.S3.AccessKey,
			AccessSecret:         config.Storage.S3.AccessSecret,
			Bucket:               config.Storage.S3.Bucket,
			Endpoint:             config.Storage.S3.Endpoint,
			InternalEndpoint:     config.Storage.S3.InternalEndpoint,
			WithInternalEndpoint: config.Storage.S3.WithInternalEndpoint,
			Region:               config.Storage.S3.Region,
			AccessMode:           AccessMode(config.Storage.S3.AccessMode),
		},
	}
}

// Instance 获取文件存储器实例
func Instance() *Storage {
	return instance(cfg)
}
