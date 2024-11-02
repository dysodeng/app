package service

import (
	"strings"
)

// ErrorMessage 错误信息类型
type ErrorMessage string

// String
func (em ErrorMessage) String() string {
	return string(em)
}

// Error
func (em ErrorMessage) Error() string {
	return em.String()
}

// Param 参数设置
func (em ErrorMessage) Param(param map[string]string) ErrorMessage {
	str := em.String()
	for k, v := range param {
		str = strings.Replace(str, "{"+k+"}", v, -1)
	}
	return ErrorMessage(str)
}

const (
	EMInputStringLengthError ErrorMessage = "{name}不能超过{length}个字，请重新输入。"
	EMMissUserIdError        ErrorMessage = "缺少用户ID"
	EMMissGoodsIdError       ErrorMessage = "缺少商品ID"
	EMMissUserTypeError      ErrorMessage = "用户类型错误"
	EMAccountDisable         ErrorMessage = "你的账号已经被系统禁止登录，请寻求客服帮助"
	EMValidCodeError         ErrorMessage = "验证码错误"
	EMValidCodeExpireError   ErrorMessage = "验证码已过期"
	EMValidCodeLimitError    ErrorMessage = "操作太频繁，请稍后再尝试"
)

const (
	EMFileSizeLimitError ErrorMessage = "当前选择的文件大小超过限制，请按照上传规定的要求处理后重新上传。"
	EMFileTypeError      ErrorMessage = "当前选择的文件格式不支持，请按照上传规定的要求处理后重新上传。"
)
