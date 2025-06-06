package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/dysodeng/fs"
	"github.com/dysodeng/fs/driver/alioss"
	"github.com/dysodeng/fs/driver/hwobs"
	"github.com/dysodeng/fs/driver/local"
	"github.com/dysodeng/fs/driver/minio"
	"github.com/dysodeng/fs/driver/s3"
	"github.com/dysodeng/fs/driver/txcos"
)

var (
	fsInstance     *Storage
	fsInstanceOnce sync.Once
)

type Storage struct {
	driver    fs.FileSystem
	cdnDomain string
}

func (storage *Storage) FileSystem() fs.FileSystem {
	return storage.driver
}

func (storage *Storage) CdnDomain() string {
	return storage.cdnDomain
}

func (storage *Storage) SignFullUrl(ctx context.Context, path string) string {
	var opts []fs.Option
	if storage.cdnDomain != "" {
		opts = append(opts, fs.WithCdnDomain(storage.cdnDomain))
	}

	fullUrl, err := storage.driver.SignFullUrl(ctx, path, opts...)
	if err != nil {
		return path
	}
	return fullUrl
}

func (storage *Storage) FullUrl(ctx context.Context, path string) string {
	var opts []fs.Option
	if storage.cdnDomain != "" {
		opts = append(opts, fs.WithCdnDomain(storage.cdnDomain))
	}

	fullUrl, err := storage.driver.FullUrl(ctx, path, opts...)
	if err != nil {
		return path
	}
	return fullUrl
}

func (storage *Storage) RelativePath(ctx context.Context, path string) string {
	originalPath, err := storage.driver.RelativePath(ctx, path)
	if err != nil {
		return path
	}
	return originalPath
}

func instance(cfg Config) *Storage {
	fsInstanceOnce.Do(func() {
		fsInstance = new(Storage)
		driver, err := generateStorageDriver(cfg)
		if err != nil {
			panic(err)
		}
		fsInstance.driver = driver
		fsInstance.cdnDomain = cfg.CdnDomain
	})
	return fsInstance
}

func generateStorageDriver(cfg Config) (fs.FileSystem, error) {
	var driver fs.FileSystem
	var err error
	switch cfg.Driver {
	case "local":
		driver, err = local.New(local.Config{
			RootPath:         cfg.Local.RootPath,
			MultipartStorage: cfg.Local.MultipartStorage,
		})
	case "minio":
		endpoint := cfg.Minio.Endpoint
		if cfg.Minio.WithInternalEndpoint {
			endpoint = cfg.Minio.InternalEndpoint
		}
		var accessMode fs.AccessMode
		switch cfg.Minio.AccessMode {
		case PublicRead:
			accessMode = fs.PublicRead
		case PublicReadWrite:
			accessMode = fs.PublicReadWrite
		default:
			accessMode = fs.Private
		}
		driver, err = minio.New(minio.Config{
			Endpoint:        endpoint,
			AccessKeyID:     cfg.Minio.AccessKey,
			SecretAccessKey: cfg.Minio.AccessSecret,
			BucketName:      cfg.Minio.Bucket,
			UseSSL:          cfg.Minio.UseSSL,
			Location:        cfg.Minio.Region,
			AccessMode:      accessMode,
		})
	case "ali_oss":
		endpoint := cfg.AliOss.Endpoint
		if cfg.AliOss.WithInternalEndpoint {
			endpoint = cfg.AliOss.InternalEndpoint
		}
		var accessMode fs.AccessMode
		switch cfg.AliOss.AccessMode {
		case PublicRead:
			accessMode = fs.PublicRead
		case PublicReadWrite:
			accessMode = fs.PublicReadWrite
		default:
			accessMode = fs.Private
		}
		driver, err = alioss.New(alioss.Config{
			Endpoint:        endpoint,
			AccessKeyID:     cfg.AliOss.AccessKey,
			SecretAccessKey: cfg.AliOss.AccessSecret,
			BucketName:      cfg.AliOss.Bucket,
			AccessMode:      accessMode,
		})
	case "hw_obs":
		endpoint := cfg.HwObs.Endpoint
		if cfg.HwObs.WithInternalEndpoint {
			endpoint = cfg.HwObs.InternalEndpoint
		}
		var accessMode fs.AccessMode
		switch cfg.HwObs.AccessMode {
		case PublicRead:
			accessMode = fs.PublicRead
		case PublicReadWrite:
			accessMode = fs.PublicReadWrite
		default:
			accessMode = fs.Private
		}
		driver, err = hwobs.New(hwobs.Config{
			Endpoint:        endpoint,
			AccessKeyID:     cfg.HwObs.AccessKey,
			SecretAccessKey: cfg.HwObs.AccessSecret,
			BucketName:      cfg.HwObs.Bucket,
			AccessMode:      accessMode,
		})
	case "tx_cos":
		endpoint := cfg.TxCos.Endpoint
		if cfg.TxCos.WithInternalEndpoint {
			endpoint = cfg.TxCos.InternalEndpoint
		}
		var accessMode fs.AccessMode
		switch cfg.TxCos.AccessMode {
		case PublicRead:
			accessMode = fs.PublicRead
		case PublicReadWrite:
			accessMode = fs.PublicReadWrite
		default:
			accessMode = fs.Private
		}
		driver, err = txcos.New(txcos.Config{
			BucketURL:  endpoint,
			SecretID:   cfg.TxCos.AccessKey,
			SecretKey:  cfg.TxCos.AccessSecret,
			AccessMode: accessMode,
		})
	case "s3":
		endpoint := cfg.S3.Endpoint
		if cfg.S3.WithInternalEndpoint {
			endpoint = cfg.S3.InternalEndpoint
		}
		var accessMode fs.AccessMode
		switch cfg.S3.AccessMode {
		case PublicRead:
			accessMode = fs.PublicRead
		case PublicReadWrite:
			accessMode = fs.PublicReadWrite
		default:
			accessMode = fs.Private
		}
		driver, err = s3.New(s3.Config{
			Endpoint:        endpoint,
			AccessKeyID:     cfg.S3.AccessKey,
			SecretAccessKey: cfg.S3.AccessSecret,
			BucketName:      cfg.S3.Bucket,
			UsePathStyle:    false,
			AccessMode:      accessMode,
		})
	default:
		return nil, fmt.Errorf("unknown storage driver: %s", cfg.Driver)
	}
	if err != nil {
		return nil, err
	}
	return driver, nil
}
