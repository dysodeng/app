package filesystem

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/dysodeng/app/internal/pkg/helper"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/filesystem/adapter"

	"github.com/pkg/errors"
)

type Filesystem struct {
	storage  adapter.Adapter
	userType string
	userId   uint64
}

type Size interface {
	Size() int64
}

type Stat interface {
	Stat() (os.FileInfo, error)
}

// Info 文件信息
type Info struct {
	Id       uint64 `json:"id"`
	Path     string `json:"path"`
	FileMd5  string `json:"file_md5"`
	Sha1     string `json:"sha1"`
	Filename string `json:"filename"`
	Ext      string `json:"ext"`
	MimeType string `json:"mime_type"`
	FileSize uint64 `json:"file_size"`
	Width    uint64 `json:"width"`
	Height   uint64 `json:"height"`
	Attach   string `json:"attach"`
	Code     int    `json:"code"`
	IsImage  uint8  `json:"is_image"`
}

// NewFilesystem 创建Filesystem
func NewFilesystem(userType string, userId uint64) (*Filesystem, error) {
	if !helper.Contain([]string{"user", "ams", "anonymous"}, userType) {
		return nil, errors.New("用户类型错误")
	}

	file := new(Filesystem)
	file.userType = userType
	file.userId = userId

	storage, err := generateAdapter()
	if err != nil {
		return nil, err
	}
	file.storage = storage

	return file, nil
}

// HasFile 判断文件是否存在
func (filesystem *Filesystem) HasFile(filePath string) bool {
	result := filesystem.storage.HasFile(filePath)
	return result
}

// DeleteFile 删除文件
func (filesystem *Filesystem) DeleteFile(filePath string) (bool, error) {
	_, err := filesystem.storage.Delete(filePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteMultipleFile 删除多个文件
func (filesystem *Filesystem) DeleteMultipleFile(filePath []string) (bool, error) {
	_, err := filesystem.storage.MultipleDelete(filePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

// FullPath 完整路径
func (filesystem *Filesystem) FullPath(filePath string) string {
	return filesystem.storage.FullPath(filePath)
}

// OriginalPath 原始文件路径
func (filesystem *Filesystem) OriginalPath(filePath string) string {
	return filesystem.storage.OriginalPath(filePath)
}

// Upload 文件上传
func (filesystem *Filesystem) Upload(file *multipart.FileHeader, allow config.FileAllow, field string) (Info, error) {
	uploader := NewUploader(filesystem.storage, allow, field)
	return uploader.Upload(filesystem.userType, file)
}

// EditorUpload 编辑器文件上传
func (filesystem *Filesystem) EditorUpload(fileBytes []byte, allow config.FileAllow, field, filename, mime, ext string) (Info, error) {
	uploader := NewUploader(filesystem.storage, allow, field)
	return uploader.EditorUpload(filesystem.userType, filesystem.userId, fileBytes, filename, mime, ext)
}

func (filesystem *Filesystem) SaveFile(dstFile string, srcFile io.Reader, mimeType string) (bool, error) {
	return filesystem.storage.Save(dstFile, srcFile, mimeType)
}

func (filesystem *Filesystem) HasDir(file string) bool {
	return filesystem.storage.HasDir(file)
}

func (filesystem *Filesystem) MkDir(dir string, mode os.FileMode) (bool, error) {
	return filesystem.storage.MkDir(dir, mode)
}

func (filesystem *Filesystem) Storage() adapter.Adapter {
	return filesystem.storage
}

// NetworkFileSaveObject 保存网络文件到存储器
func (filesystem *Filesystem) NetworkFileSaveObject(networkFileUrl, filePath, mime string) (bool, error) {
	res, err := http.Get(networkFileUrl)
	if err != nil {
		return false, err
	}
	defer func() {
		res.Body.Close()
	}()

	reader := bufio.NewReaderSize(res.Body, 32*1024)

	return filesystem.storage.Save(filePath, reader, mime)
}

// DownloadNetworkFile 下载网络文件到本地
func (filesystem *Filesystem) DownloadNetworkFile(networkFileUrl, filePath string) error {
	response, err := http.Get(networkFileUrl)
	if err != nil {
		fmt.Println("Error while downloading", networkFileUrl, ":", err)
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading response body:", err)
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("Error while saving file:", err)
		return err
	}

	return nil
}

// FileExists 查询文件是否存在
func FileExists(userType, userId string, sha1 string, md5 string) (Info, error) {
	switch userType {
	default:
		return Info{}, errors.New("用户类型错误")
	}
}

// DeleteFile 删除不存在的文件记录
func DeleteFile(userType string, id uint64) {
	switch userType {
	}
}

// SaveFile 保存文件
func SaveFile(userType string, userId uint64, info Info) (id uint64, err error) {
	switch userType {
	case "user":
		return 0, nil
	case "ams":
		return 0, nil
	case "anonymous":
		return 0, nil
	}
	return 0, errors.New("用户类型错误")
}

var filesystemInstance *Filesystem
var filesystemInstanceOnce sync.Once

// Instance 获取单例
// 注意，单例仅仅为了获取filesystem及storage实例，不涉及任何上传相关的实例，如果要上传文件，还是得使用NewFilesystem生成实例
func Instance() *Filesystem {

	filesystemInstanceOnce.Do(func() {

		filesystemInstance = new(Filesystem)

		storage, err := generateAdapter()
		if err != nil {
			panic("file storage error:" + err.Error())
		}
		filesystemInstance.storage = storage
	})

	return filesystemInstance
}

func generateAdapter() (adapter.Adapter, error) {
	var storage adapter.Adapter
	switch config.Filesystem.Storage {
	case "ali_oss": // 阿里云OSS
		storage = adapter.NewAliOssAdapter(adapter.AliOssConfig{
			AccessId:   config.Filesystem.AliOss.AccessKey,
			AccessKey:  config.Filesystem.AliOss.AccessSecret,
			EndPoint:   config.Filesystem.AliOss.Endpoint,
			Region:     config.Filesystem.AliOss.Region,
			BucketName: config.Filesystem.AliOss.Bucket,
			StsRoleArn: config.Filesystem.AliOss.StsRoleArn,
		})

	case "hw_obs": // 华为云OBS
		storage = adapter.NewHwObsAdapter(adapter.HwObsConfig{
			AccessKey:  config.Filesystem.HwObs.AccessKey,
			SecretKey:  config.Filesystem.HwObs.AccessSecret,
			EndPoint:   config.Filesystem.HwObs.Endpoint,
			BucketName: config.Filesystem.HwObs.Bucket,
		})

	case "minio": // MinIO
		storage = adapter.NewMinioAdapter(adapter.MinioConfig{
			AccessKey:  config.Filesystem.Minio.AccessKey,
			SecretKey:  config.Filesystem.Minio.AccessSecret,
			EndPoint:   config.Filesystem.Minio.Endpoint,
			BucketName: config.Filesystem.Minio.Bucket,
			UseSSL:     config.Filesystem.Minio.UseSSL,
		})

	case "local": // 本地文件系统
		storage = adapter.NewLocalAdapter(adapter.LocalConfig{
			BasePath:  config.Filesystem.Local.BasePath,
			LogicPath: config.Filesystem.Local.LogicPath,
			BaseUrl:   config.Filesystem.Local.BaseUrl,
		})

	default:
		log.Println("file storage error:" + config.Filesystem.Storage + " not found.")
		return nil, errors.New("文件存储驱动错误")
	}

	return storage, nil
}
