package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// Md5 信息摘要
func Md5(content string) string {
	data := []byte(content)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

// Sha256 sha256信息摘要
func Sha256(content []byte) ([]byte, error) {
	h := sha256.New()
	n, err := h.Write(content)
	if err != nil {
		return nil, err
	}

	if n != len(content) {
		return nil, errors.New("write length error")
	}

	return h.Sum(nil), nil
}

// FileMd5 获取文件md5值
func FileMd5(filename string) (fileMd5 string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	// 计算文件md5
	fm := md5.New()
	_, _ = io.Copy(fm, file)
	fileMd5 = hex.EncodeToString(fm.Sum(nil))

	return
}

// FileSha1 获取文件sha1值
func FileSha1(filename string) (fileSha1 string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	// 计算文件sha1
	fSha1 := sha1.New()
	_, _ = io.Copy(fSha1, file)
	fileSha1 = hex.EncodeToString(fSha1.Sum(nil))

	return
}
