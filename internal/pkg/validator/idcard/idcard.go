package idcard

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// 校验系数
var cardCoefficient = []int64{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// 校验码
var checksum = []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

// Check 身份证有效性验证
// @param idCard string 身份证号
func Check(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}
	result, _ := regexp.MatchString(`^\d{17}[\dxX]$`, idCard)
	if !result {
		return false
	}

	pre17 := idCard[:17]                   // 前17位
	suffix := strings.ToUpper(idCard[17:]) // 第18位为检验码

	length := len(pre17)

	// 前17位与校验系数相乘
	var sum int64 = 0
	for i := 0; i < length; i++ {
		currentNum, err := strconv.ParseInt(string(pre17[i]), 10, 64)
		if err != nil {
			return false
		}
		sum += cardCoefficient[i] * currentNum
	}

	// 取模
	seek := sum % 11

	if suffix == checksum[seek] {
		return true
	}

	return false
}

// Hide 隐藏身份证号中间位数
// @param idCard string 身份证号
// @param symbol string 隐藏标示符号
func Hide(idCard, symbol string) string {
	if symbol == "" {
		symbol = "*"
	}

	symbol = strings.Repeat(symbol, 8)

	newIdCard := idCard
	if Check(idCard) {
		hideStr := idCard[4:14]
		newIdCard = strings.Replace(idCard, hideStr, symbol, -1)
	}

	return newIdCard
}

// HideRealName 隐藏姓名
// @param realName string 姓名
// @param symbol string 隐藏标示符号
func HideRealName(realName, symbol string) string {
	if symbol == "" {
		symbol = "*"
	}
	length := utf8.RuneCountInString(realName)
	if length > 1 {
		name := []rune(realName)
		suffix := string(name[length-1])
		return strings.Repeat(symbol, length-1) + suffix
	}
	return realName
}

// weightedFactors 信用代码加权因子
var weightedFactors = []int64{1, 3, 9, 27, 19, 26, 16, 17, 20, 29, 25, 13, 8, 24, 10, 30, 28}

// availableString 可用字符，去掉了I,O,S,V,Z等字符
var availableString = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "T", "U", "W", "X", "Y"}

// CompanyCreditCodeCheck 企业统一社会信用代码有效性验证
func CompanyCreditCodeCheck(creditCode string) bool {
	length := len(creditCode)
	if length != 18 {
		return false
	}
	result, _ := regexp.MatchString(`^[\dA-Z]+$`, creditCode)
	if !result {
		return false
	}

	var sum int64 = 0

	for i := 0; i < length-1; i++ {
		anCode := creditCode[i]
		if pos := findFirstRange(availableString, string(anCode)); pos != -1 {
			sum += int64(pos) * weightedFactors[i] // 权重与加权因子相乘之和
		} else {
			return false
		}
	}

	seek := 31 - sum%31
	if seek == 31 {
		seek = 0
	}

	logicCode := availableString[seek]
	checkCode := creditCode[17:]

	return logicCode == checkCode
}

// CompanyCreditCodeHide 隐藏企业统一社会信用代码中间位数
func CompanyCreditCodeHide(creditCode, symbol string) string {
	if symbol == "" {
		symbol = "*"
	}

	symbol = strings.Repeat(symbol, 14)

	newCreditCode := creditCode
	if CompanyCreditCodeCheck(creditCode) {
		hideStr := creditCode[1:17]
		newCreditCode = strings.Replace(creditCode, hideStr, symbol, -1)
	}

	return newCreditCode
}

// CompanyNameHide 隐藏企业名称
// @param companyName string 企业名称
// @param symbol string 隐藏标示符号
func CompanyNameHide(companyName, symbol string) string {
	if symbol == "" {
		symbol = "*"
	}
	length := utf8.RuneCountInString(companyName)
	if length > 2 {
		name := []rune(companyName)
		var prefix string
		var suffix string
		if length > 4 {
			prefix = string(name[0:2])
			suffix = string(name[length-2:])
		} else {
			prefix = string(name[0])
			suffix = string(name[length-1])
		}
		return prefix + strings.Repeat(symbol, 14) + suffix
	}
	return companyName
}

func findFirstRange(str []string, target string) int {
	var pos = -1
	if len(str) <= 0 {
		return pos
	}

	length := len(str)
	for i := 0; i < length; i++ {
		if str[i] == target {
			pos = i
			break
		}
	}

	return pos
}
