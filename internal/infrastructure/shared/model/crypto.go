package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"strings"

	"github.com/dysodeng/app/internal/infrastructure/shared/crypto/aes"
	"github.com/dysodeng/app/internal/infrastructure/shared/helper"
)

// Crypto 加密字段
type Crypto string

func (c Crypto) Value() (driver.Value, error) {
	key := helper.RandomString(16, helper.ModeAlphanumeric)
	iv := helper.RandomString(16, helper.ModeAlphanumeric)

	crypto, _ := aes.Encrypt(
		helper.StringToBytes(c),
		helper.StringToBytes(key),
		helper.StringToBytes(iv),
	)

	ivKey := key + "." + iv
	cryptoBytes := make([]byte, base64.StdEncoding.EncodedLen(len(crypto)))
	base64.StdEncoding.Encode(cryptoBytes, crypto)
	ivKeyBytes := make([]byte, base64.StdEncoding.EncodedLen(len(ivKey)))
	base64.StdEncoding.Encode(ivKeyBytes, helper.StringToBytes(ivKey))

	buf := bytes.NewBuffer(nil)
	buf.Write(cryptoBytes)
	buf.WriteString(".")
	buf.Write(ivKeyBytes)

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (c *Crypto) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	orgBase64Value := helper.BytesToString(v.([]byte))
	if orgBase64Value == "" {
		return nil
	}

	base64Value, err := base64.StdEncoding.DecodeString(orgBase64Value)
	if err != nil {
		return nil
	}

	array := strings.Split(helper.BytesToString(base64Value), ".")
	if len(array) != 2 {
		return nil
	}
	base64Key, err := base64.StdEncoding.DecodeString(array[1])
	if err != nil {
		return nil
	}
	k := strings.Split(helper.BytesToString(base64Key), ".")
	if len(k) != 2 {
		return nil
	}

	key := k[0]
	iv := k[1]

	b, err := base64.StdEncoding.DecodeString(array[0])
	if err != nil {
		return nil
	}

	descByte, err := aes.Decrypt(b, helper.StringToBytes(key), helper.StringToBytes(iv))
	if err != nil {
		return nil
	}
	*c = Crypto(descByte)
	return nil
}

func (c Crypto) String() string {
	return string(c)
}
