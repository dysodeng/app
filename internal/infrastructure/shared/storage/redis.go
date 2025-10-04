package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dysodeng/fs/driver/local"
	"github.com/redis/go-redis/v9"
)

// RedisMultipartStorage Redis实现的本地文件分片上传状态存储
type RedisMultipartStorage struct {
	client redis.Cmdable
	prefix string // Redis key前缀
}

func NewRedisMultipartStorage(client redis.Cmdable, prefix string) *RedisMultipartStorage {
	if prefix == "" {
		prefix = "fs:multipart:"
	}
	return &RedisMultipartStorage{client: client, prefix: prefix}
}

func (s *RedisMultipartStorage) getKey(uploadID string) string {
	return s.prefix + uploadID
}

func (s *RedisMultipartStorage) Save(upload *local.MultipartUpload) error {
	data, err := json.Marshal(*upload)
	if err != nil {
		return err
	}
	return s.client.Set(context.Background(), s.getKey(upload.UploadID), data, 0).Err()
}

func (s *RedisMultipartStorage) Get(uploadID string) (*local.MultipartUpload, error) {
	data, err := s.client.Get(context.Background(), s.getKey(uploadID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("upload ID not found")
		}
		return nil, err
	}

	upload := &local.MultipartUpload{}
	if err := json.Unmarshal(data, upload); err != nil {
		return nil, err
	}
	return upload, nil
}

func (s *RedisMultipartStorage) Delete(uploadID string) error {
	return s.client.Del(context.Background(), s.getKey(uploadID)).Err()
}

func (s *RedisMultipartStorage) List() ([]*local.MultipartUpload, error) {
	keys, err := s.client.Keys(context.Background(), s.prefix+"*").Result()
	if err != nil {
		return nil, err
	}

	var uploads []*local.MultipartUpload
	for _, key := range keys {
		uploadID := key[len(s.prefix):]
		if upload, err := s.Get(uploadID); err == nil {
			uploads = append(uploads, upload)
		}
	}
	return uploads, nil
}
