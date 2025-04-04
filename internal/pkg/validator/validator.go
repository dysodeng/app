package validator

import (
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/dysodeng/app/internal/pkg/validator/idcard"
)

// IsMobile 判断是否为手机号
func IsMobile(value string) bool {
	result, _ := regexp.MatchString(`^(1[0-9][0-9]\d{4,8})$`, value)
	return result
}

// IsPhone 判断是否为固定电话号码
func IsPhone(value string) bool {
	result, _ := regexp.MatchString(`^(\d{4}-|\d{3}-)?(\d{8}|\d{7})$`, value)
	return result
}

// IsPhone400 判断是否为400电话
func IsPhone400(value string) bool {
	result, _ := regexp.MatchString(`^400(-\d{3,4}){2}$`, value)
	return result
}

// IsTel 判断是否为电话号码
func IsTel(value string) bool {
	return !(!IsMobile(value) && !IsPhone(value) && !IsPhone400(value))
}

// IsEmail 判断是否为邮箱
func IsEmail(value string) bool {
	result, _ := regexp.MatchString(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, value)
	return result
}

// IsUrl 判断是否为网址
func IsUrl(value string) bool {
	result, _ := regexp.MatchString(`^(http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?$`, value)
	return result
}

// IsSSLUrl 判断是否为https链接
func IsSSLUrl(value string) bool {
	result, _ := regexp.MatchString(`^https:\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?$`, value)
	return result
}

// IsIdCard 是否身份证号
func IsIdCard(value string) bool {
	return idcard.Check(value)
}

// IsCardAdult 验证身份证号是否为成年人
func IsCardAdult(idCard string) bool {
	if !IsIdCard(idCard) {
		return false
	}

	yearString := idCard[6:10]
	year, _ := strconv.ParseUint(yearString, 10, 64)

	return uint64(time.Now().Year())-year >= 18
}

// IsBankCard 是否银行卡号
func IsBankCard(value string) bool {
	result, _ := regexp.MatchString(`^(\d{16}|\d{17}|\d{18}|\d{19})$`, value)
	return result
}

// IsDate 是否为日期(年-月-日)
func IsDate(value string) bool {
	result, _ := regexp.MatchString(`^(\d{4}|\d{2})-((0?([1-9]))|(1[0-2]))-((0?[1-9])|([12]([0-9]))|(3[0|1]))$`, value)
	return result
}

// IsShortDate 是否为日期(年-月)
func IsShortDate(value string) bool {
	result, _ := regexp.MatchString(`^(\d{4}|\d{2})-((0?([1-9]))|(1[0-2]))$`, value)
	return result
}

// IsUintNumber 是否为整型数字
func IsUintNumber(value string) bool {
	if value == "" {
		return false
	}

	if value == "0" {
		return true
	} else {
		number, _ := strconv.ParseUint(value, 10, 64)
		if number == 0 {
			return false
		}
	}

	return true
}

// IsSafePassword 是否为安全的密码
func IsSafePassword(value string, length int) bool {
	if length <= 0 {
		length = 6
	}
	// 是否为纯字符
	result, _ := regexp.MatchString(`^([a-zA-Z])+$`, value)
	if result {
		return false
	}
	// 是否为纯数字
	result, _ = regexp.MatchString(`^\d+$`, value)
	if result {
		return false
	}
	// 是否为纯符号
	result, _ = regexp.MatchString(`^\W+$`, value)
	if result {
		return false
	}
	if len(value) < length {
		return false
	}
	return true
}

// IsDomain 是否为域名
func IsDomain(domain string) bool {
	switch {
	case len(domain) == 0:
		return false
	case len(domain) > 255:
		return false
	}

	var l int
	for i := 0; i < len(domain); i++ {
		b := domain[i]
		if b == '.' {
			// check domain labels validity
			switch {
			case i == l:
				return false
			case i-l > 63:
				return false
			case domain[l] == '-':
				return false
			case domain[i-1] == '-':
				return false
			}
			l = i + 1
			continue
		}
		// test label character validity, note: tests are ordered by decreasing validity frequency
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b >= 'A' && b <= 'Z') {
			// show the printable unicode character starting at byte offset i
			c, _ := utf8.DecodeRuneInString(domain[i:])
			if c == utf8.RuneError {
				return false
			}
			return false
		}
	}

	// check top level domain validity
	switch {
	case l == len(domain):
		return false
	case len(domain)-l > 63:
		return false
	case domain[l] == '-':
		return false
	case domain[len(domain)-1] == '-':
		return false
	case domain[l] >= '0' && domain[l] <= '9':
		return false
	}

	return true
}
