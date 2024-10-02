package model

import (
	"database/sql/driver"
	"encoding/base64"
	"strings"

	"github.com/dysodeng/app/internal/pkg/helper"

	"github.com/dysodeng/app/internal/pkg/crypto/aes"
)

// Crypto 加密字段
type Crypto string

func (c Crypto) Value() (driver.Value, error) {
	key := helper.RandomStringBytesMask(16)
	iv := helper.RandomStringBytesMask(16)
	crypto, _ := aes.Encrypt([]byte(c), []byte(key), []byte(iv))
	return base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString(crypto) + "." + base64.StdEncoding.EncodeToString([]byte(key+"&"+iv)))), nil
}

func (c *Crypto) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	orgBase64Value := string(v.([]byte))
	if orgBase64Value == "" {
		return nil
	}

	base64Value, err := base64.StdEncoding.DecodeString(orgBase64Value)
	if err != nil {
		return nil
	}

	array := strings.Split(string(base64Value), ".")
	if len(array) != 2 {
		return nil
	}
	base64Key, err := base64.StdEncoding.DecodeString(array[1])
	if err != nil {
		return nil
	}
	k := strings.Split(string(base64Key), "&")
	if len(k) != 2 {
		return nil
	}

	key := k[0]
	iv := k[1]

	b, err := base64.StdEncoding.DecodeString(array[0])
	if err != nil {
		return nil
	}

	descByte, err := aes.Decrypt(b, []byte(key), []byte(iv))
	if err != nil {
		return nil
	}

	*c = Crypto(descByte)

	return nil
}

func (c Crypto) String() string {
	return string(c)
}
