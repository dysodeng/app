package qrcode

import (
	"bytes"
	"fmt"
	"image"
	"net/url"

	"github.com/dysodeng/app/internal/config"

	"github.com/yeqown/go-qrcode"

	"github.com/pkg/errors"
)

type QrCode struct {
	content string
	logo    image.Image
	ptr     *qrcode.QRCode
}

// NewQrCode 创建二维码
func NewQrCode(text string, border ...int) (*QrCode, error) {
	opts := []qrcode.ImageOption{
		qrcode.WithBuiltinImageEncoder(qrcode.PNG_FORMAT),
	}
	if len(border) > 0 {
		opts = append(opts, qrcode.WithBorderWidth(border...))
	}

	ptr, err := qrcode.New(text, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create QR Code error")
	}

	return &QrCode{
		content: text,
		ptr:     ptr,
	}, nil
}

// Text 获取二维码内容
func (qr *QrCode) Text() string {
	return qr.content
}

// SaveToBuffer 保存到buffer
func (qr *QrCode) SaveToBuffer() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	err := qr.ptr.SaveTo(buf)
	if err != nil {
		return nil, errors.Wrap(err, "save to buffer QR Code error")
	}
	return buf, nil
}

// SaveToFile 保存到文件
func (qr *QrCode) SaveToFile(saveToPath string) error {
	err := qr.ptr.Save(saveToPath)
	if err != nil {
		return errors.Wrap(err, "save to file QR Code error")
	}
	return nil
}

// Url 生成二维码链接
func Url(text string, isUrl bool) string {
	textIsUrl := "0"
	if isUrl {
		textIsUrl = "1"
		text = url.QueryEscape(text)
	}
	return fmt.Sprintf("%s/api/v1/common/qr_code?is_url=%s&text=%s", config.App.Domain, textIsUrl, text)
}
