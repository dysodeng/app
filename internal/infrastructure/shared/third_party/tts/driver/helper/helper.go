package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"path"
	"strings"
	"time"
)

func GenerateFilePath(ext string) (string, error) {
	if strings.ContainsRune(ext, '/') { // 防止路径注入
		ext = ".invalid"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ext = "." + strings.ReplaceAll(ext, ".", "")

	now := time.Now()
	dateDir := now.Format("2006/01/02")

	// 生成唯一文件名部分（纳秒时间戳+随机熵）
	randomBytes := make([]byte, 4) // 4字节提供32位熵值
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// 构建最终文件名（格式：时间戳_随机数.扩展名）
	fileName := fmt.Sprintf("%d_%s%s",
		now.UnixNano(),
		base64.RawURLEncoding.EncodeToString(randomBytes),
		ext,
	)

	return path.Join(
		"resources/audio",
		dateDir,
		fileName,
	), nil
}
