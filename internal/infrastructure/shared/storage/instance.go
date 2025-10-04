package storage

import (
	"github.com/dysodeng/fs/driver/local"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
)

var cfg Config

func Init(c *config.Config) (*Storage, error) {
	var localMultipartStorage local.MultipartStorage
	if c.Storage.Driver == "local" {
		switch c.Storage.Local.MultipartStorage {
		case "redis":
			localMultipartStorage = NewRedisMultipartStorage(redis.MainClient(), "")
		default:
			// 默认为文件
			localMultipartStorage = nil
		}
	}

	cfg = Config{
		Driver:    c.Storage.Driver,
		CdnDomain: c.Storage.CdnDomain,
		Local: Local{
			RootPath:         c.Storage.Local.RootPath,
			MultipartStorage: localMultipartStorage,
		},
		Minio: CloudStorage{
			AccessKey:            c.Storage.MinIO.AccessKey,
			AccessSecret:         c.Storage.MinIO.AccessSecret,
			Bucket:               c.Storage.MinIO.Bucket,
			Endpoint:             c.Storage.MinIO.Endpoint,
			InternalEndpoint:     c.Storage.MinIO.InternalEndpoint,
			WithInternalEndpoint: c.Storage.MinIO.WithInternalEndpoint,
			Region:               c.Storage.MinIO.Region,
			UseSSL:               c.Storage.MinIO.UseSSL,
			AccessMode:           AccessMode(c.Storage.MinIO.AccessMode),
		},
		HwObs: CloudStorage{
			AccessKey:            c.Storage.HwObs.AccessKey,
			AccessSecret:         c.Storage.HwObs.AccessSecret,
			Bucket:               c.Storage.HwObs.Bucket,
			Endpoint:             c.Storage.HwObs.Endpoint,
			InternalEndpoint:     c.Storage.HwObs.InternalEndpoint,
			WithInternalEndpoint: c.Storage.HwObs.WithInternalEndpoint,
			Region:               c.Storage.HwObs.Region,
			UseSSL:               c.Storage.HwObs.UseSSL,
			AccessMode:           AccessMode(c.Storage.HwObs.AccessMode),
		},
		AliOss: CloudStorage{
			AccessKey:            c.Storage.AliOss.AccessKey,
			AccessSecret:         c.Storage.AliOss.AccessSecret,
			Bucket:               c.Storage.AliOss.Bucket,
			Endpoint:             c.Storage.AliOss.Endpoint,
			InternalEndpoint:     c.Storage.AliOss.InternalEndpoint,
			WithInternalEndpoint: c.Storage.AliOss.WithInternalEndpoint,
			Region:               c.Storage.AliOss.Region,
			UseSSL:               c.Storage.AliOss.UseSSL,
			AccessMode:           AccessMode(c.Storage.AliOss.AccessMode),
		},
		TxCos: CloudStorage{
			AccessKey:            c.Storage.TxCos.AccessKey,
			AccessSecret:         c.Storage.TxCos.AccessSecret,
			Bucket:               c.Storage.TxCos.Bucket,
			Endpoint:             c.Storage.TxCos.Endpoint,
			InternalEndpoint:     c.Storage.TxCos.InternalEndpoint,
			WithInternalEndpoint: c.Storage.TxCos.WithInternalEndpoint,
			Region:               c.Storage.TxCos.Region,
			UseSSL:               c.Storage.TxCos.UseSSL,
			AccessMode:           AccessMode(c.Storage.TxCos.AccessMode),
		},
		S3: CloudStorage{
			AccessKey:            c.Storage.S3.AccessKey,
			AccessSecret:         c.Storage.S3.AccessSecret,
			Bucket:               c.Storage.S3.Bucket,
			Endpoint:             c.Storage.S3.Endpoint,
			InternalEndpoint:     c.Storage.S3.InternalEndpoint,
			WithInternalEndpoint: c.Storage.S3.WithInternalEndpoint,
			Region:               c.Storage.S3.Region,
			AccessMode:           AccessMode(c.Storage.S3.AccessMode),
		},
	}

	return Instance(), nil
}

// Instance 获取文件存储器实例
func Instance() *Storage {
	return instance(cfg)
}
