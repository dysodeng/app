package filesystem

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/dysodeng/app/internal/pkg/helper"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/api"

	"github.com/dysodeng/filesystem/adapter"

	"github.com/pkg/errors"
)

type Uploader struct {
	storage adapter.Adapter
	allow   config.FileAllow
	field   string
}

func NewUploader(storage adapter.Adapter, allow config.FileAllow, field string) *Uploader {
	return &Uploader{storage: storage, allow: allow, field: field}
}

// Upload 文件上传
func (uploader *Uploader) Upload(userType string, fileHeader *multipart.FileHeader) (Info, error) {

	rootPath := ""

	switch userType {
	case "user": // 终端用户
		rootPath += "user/"
		break
	case "ams": // 管理端
		rootPath += "ams/file/"
		break
	case "anonymous": // 匿名用户
		rootPath += "anonymous/file/"
		break
	default:
		return Info{}, errors.New("用户类型错误")
	}

	file, err := fileHeader.Open()
	dstFileReader, err := fileHeader.Open()
	if err != nil {
		log.Println(err.Error())
		return Info{}, errors.New("文件读取错误")
	}

	// 类型与后缀
	var mime = fileHeader.Header.Get("Content-Type")
	var ext string
	filename := fileHeader.Filename

	fType, ok := MimeType[mime]
	if fType != "" {
		fType = strings.ToLower(fType)
	}

	log.Printf("upload image: mime:%s ext:%s fType:%s, filename:%s", mime, ext, fType, filename)

	if !ok || !IsExistsMimeAllow(fType, uploader.allow.AllowMimeType) {
		return Info{}, api.EMFileTypeError
	}

	extSlice := strings.Split(filename, ".")
	if len(extSlice) >= 2 {
		ext = extSlice[len(extSlice)-1]
	}

	// 计算文件大小
	var size int64
	if fileSize, ok := file.(Size); ok {
		size = fileSize.Size()
	}

	if size > uploader.allow.AllowCapacitySize {
		return Info{}, api.EMFileSizeLimitError
	}

	// 如果是图片，获取图片尺寸
	img, _, err := image.DecodeConfig(file)
	var isImage uint8
	var imageWidth, imageHeight uint64
	if err == nil {
		isImage = 1
		imageWidth = uint64(img.Width)
		imageHeight = uint64(img.Height)
	} else {
		isImage = 0
	}

	// 计算文件md5
	fileMd5 := md5.New()
	_, _ = io.Copy(fileMd5, file)
	md5String := hex.EncodeToString(fileMd5.Sum(nil))

	// 计算文件sha1
	fileSha1 := sha1.New()
	_, _ = io.Copy(fileSha1, file)
	sha1String := hex.EncodeToString(fileSha1.Sum(nil))

	savePath := time.Now().Format(time.DateOnly) + "/"
	filePath := userType + helper.CreateOrderNo()

	dstFile := rootPath + savePath + filePath //原文件
	if ext != "" {
		dstFile += "." + ext
	}

	// 创建目录
	if !uploader.storage.HasDir(rootPath + savePath) {
		_, err = uploader.storage.MkDir(rootPath+savePath, 0755)
		if err != nil {
			log.Println(err)
			return Info{}, err
		}
	}

	// 上传文件
	result, err := uploader.storage.Save(dstFile, dstFileReader, mime)
	if err != nil {
		log.Println(err.Error(), result)
		return Info{}, err
	}

	info := Info{
		Path:     dstFile,
		FileMd5:  md5String,
		Sha1:     sha1String,
		Filename: strings.Replace(filename, "."+ext, "", 1),
		Ext:      "." + ext,
		MimeType: mime,
		Width:    imageWidth,
		Height:   imageHeight,
		FileSize: uint64(size),
		IsImage:  isImage,
	}
	var filenameRune = []rune(info.Filename)
	if len(filenameRune) > 40 {
		filenameRune = filenameRune[len(filenameRune)-40:]
		info.Filename = string(filenameRune)
	}

	// 保存文件
	// id, err := SaveFile(userType, userId, info)
	//if err == nil {
	//	info.Id = id
	//}

	return info, nil
}

// EditorUpload 编辑器文件上传
func (uploader *Uploader) EditorUpload(userType string, userId uint64, fileBytes []byte, filename, mime, ext string) (Info, error) {
	if userType != "ams" {
		if userId <= 0 {
			return Info{}, errors.New("用户ID为空")
		}
	}

	rootPath := ""

	switch userType {
	case "user": // 终端用户
		if userId <= 0 {
			return Info{}, api.EMMissUserIdError
		}
		rootPath += fmt.Sprintf("user/%d/", userId)
		break
	case "ams":
		rootPath += "ams/file/"
		break
	default:
		return Info{}, errors.New("用户类型错误")
	}

	srcFile := bytes.NewReader(fileBytes)
	distFile := bytes.NewReader(fileBytes)

	fType, ok := MimeType[mime]
	if fType != "" {
		fType = strings.ToLower(fType)
	}

	log.Printf("upload image: mime:%s ext:%s fType:%s, filename:%s", mime, ext, fType, filename)

	if !ok || !IsExistsMimeAllow(fType, uploader.allow.AllowMimeType) {
		return Info{}, api.EMFileTypeError
	}

	// 计算文件大小
	var size = srcFile.Size()

	// 如果是图片，获取图片尺寸
	img, _, err := image.DecodeConfig(srcFile)
	var isImage uint8
	var imageWidth, imageHeight uint64
	if err == nil {
		isImage = 1
		imageWidth = uint64(img.Width)
		imageHeight = uint64(img.Height)
	} else {
		isImage = 0
	}

	// 计算文件md5
	fileMd5 := md5.New()
	_, _ = io.Copy(fileMd5, srcFile)
	md5String := hex.EncodeToString(fileMd5.Sum(nil))

	// 计算文件sha1
	fileSha1 := sha1.New()
	_, _ = io.Copy(fileSha1, srcFile)
	sha1String := hex.EncodeToString(fileSha1.Sum(nil))

	savePath := time.Now().Format(time.DateOnly) + "/"
	filePath := userType + helper.CreateOrderNo()

	dstFile := rootPath + savePath + filePath
	if ext != "" {
		dstFile += "." + ext
	}

	// 创建目录
	if !uploader.storage.HasDir(rootPath + savePath) {
		_, err = uploader.storage.MkDir(rootPath+savePath, 0755)
		if err != nil {
			log.Println(err)
			return Info{}, err
		}
	}

	// 上传文件
	result, err := uploader.storage.Save(dstFile, distFile, mime)
	if err != nil {
		log.Println(err.Error(), result)
		return Info{}, err
	}

	info := Info{
		Path:     dstFile,
		FileMd5:  md5String,
		Sha1:     sha1String,
		Filename: strings.Replace(filename, "."+ext, "", 1),
		Ext:      "." + ext,
		MimeType: mime,
		Width:    imageWidth,
		Height:   imageHeight,
		FileSize: uint64(size),
		IsImage:  isImage,
	}
	var filenameRune = []rune(info.Filename)
	if len(filenameRune) > 40 {
		filenameRune = filenameRune[len(filenameRune)-40:]
		info.Filename = string(filenameRune)
	}

	// 保存文件
	id, err := SaveFile(userType, userId, info)
	if err == nil {
		info.Id = id
	}

	return info, nil
}

func (uploader *Uploader) HasFile(filePath string) bool {
	return uploader.storage.HasFile(filePath)
}
