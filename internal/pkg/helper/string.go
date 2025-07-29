package helper

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// RandomStringMode 随机字符串模式
type RandomStringMode int

const (
	ModeNumber       RandomStringMode = iota // 纯数字
	ModeLetter                               // 纯字母
	ModeAlphanumeric                         // 字母和数字
	ModeComplex                              // 复杂
)

// GeneratePassword 生成密码
// @param string password 明文密码
// @return string
func GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(StringToBytes(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return BytesToString(hash), nil
}

// VerifyPassword 验证密码
// @param string hashedPassword hash密码
// @param string plainPassword 明文密码
func VerifyPassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(
		StringToBytes(hashedPassword),
		StringToBytes(plainPassword),
	)
}

// RandomString 生成随机字符串
// @param int length 生成字符串长度
// @return string
func RandomString(length int, mode RandomStringMode) string {
	switch mode {
	case ModeNumber:
		return lo.RandomString(length, lo.NumbersCharset)
	case ModeLetter:
		return lo.RandomString(length, lo.LettersCharset)
	case ModeAlphanumeric:
		return lo.RandomString(length, lo.AlphanumericCharset)
	case ModeComplex:
		fallthrough
	default:
		charset := append(lo.AlphanumericCharset, []rune("!@#$%^&*()_+-=[],./;<>?")...)
		return lo.RandomString(length, charset)
	}
}

// RandomNumberString 生成随机数字字符串
// @param int length 生成字符串长度
func RandomNumberString(length int) string {
	return RandomString(length, ModeNumber)
}

// GbkToUtf8 GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk Utf8 转 gbk
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// CreateOrderNo 生成唯一订单号
// @return string
func CreateOrderNo() string {
	builder := strings.Builder{}
	builder.WriteString(time.Now().Format("20060102150405"))

	t := time.Now().UnixNano()
	s := strconv.FormatInt(t, 10)
	b := BytesToString(StringToBytes(s)[len(s)-9:])
	c := BytesToString(StringToBytes(b)[:7])

	builder.WriteString(c)
	builder.WriteString(RandomNumberString(6))
	return builder.String()
}

// ReplaceString 批量替换字符串
func ReplaceString(str string, findSlice, replaceSlice []string) string {
	if len(findSlice) != len(replaceSlice) {
		return str
	}

	for i, find := range findSlice {
		str = strings.ReplaceAll(str, find, replaceSlice[i])
	}

	return str
}

// HideCellphone 隐藏手机号码中间四位
func HideCellphone(cellphone string) string {
	if len(cellphone) > 7 {
		return cellphone[:3] + "****" + cellphone[7:]
	} else if len(cellphone) > 3 {
		return cellphone[:3] + "****"
	} else if len(cellphone) > 0 {
		return "****"
	}
	return cellphone
}

func HideEmail(email string) string {
	if email == "" {
		return ""
	}
	temp := strings.Split(email, "@")
	if len(temp) != 2 || len(temp[0]) == 0 || len(temp[1]) == 0 {
		return email[0:1] + "****"
	}
	prefix := temp[0]
	suffix := temp[1]
	return prefix[0:1] + "****@" + suffix
}

// HideRealName 隐藏姓名中间文字
func HideRealName(name string) string {
	n := []rune(name)
	l := len(n)
	var newName string
	if l >= 2 {
		if l > 2 {
			newName = string(n[0]) + "*" + string(n[l-1])
		} else {
			newName = string(n[0]) + "*"
		}
	} else {
		return name
	}
	return newName
}

// StringToBytes 字符串转字节数组
func StringToBytes[T ~string](s T) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0], x[1], x[1]}
	res := *(*[]byte)(unsafe.Pointer(&b))
	return res
}

// BytesToString 字节数组转字符串
func BytesToString[T ~[]byte](b T) string {
	return *(*string)(unsafe.Pointer(&b))
}
