package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	pkgCrypto "github.com/dysodeng/app/internal/pkg/crypto"

	"github.com/pkg/errors"
)

// Encrypt 签名
// @param content string 原始内容
// @param privateKey string 加密私钥
func Encrypt(content, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private_key error")
	}

	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPrivate, _ := private.(*rsa.PrivateKey)

	hashed, err := pkgCrypto.Sha256([]byte(content))
	if err != nil {
		return "", err
	}

	sign, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivate, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// Check 验证签名
// @param content string 待验签内容
// @param sign string 签名串
// @param publicKey string 公钥
func Check(content, sign, publicKey string) (bool, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, errors.New("public_key error")
	}

	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	rsaPublic, _ := public.(*rsa.PublicKey)

	digest, err := pkgCrypto.Sha256([]byte(content))
	if err != nil {
		return false, err
	}

	data, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(rsaPublic, crypto.SHA256, digest, data)
	if err != nil {
		return false, err
	}

	return true, nil
}
