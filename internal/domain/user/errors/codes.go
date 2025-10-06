package errors

import (
	domainErrors "github.com/dysodeng/app/internal/domain/shared/errors"
)

// 用户领域错误码前缀
const (
	// PrefixUser 用户领域错误码前缀
	PrefixUser = "USER"
)

// 用户领域错误码
const (
	CodeUserNotFound               = "USER_NOT_FOUND"
	CodeUserAlreadyExists          = "USER_ALREADY_EXISTS"
	CodeUserDisabled               = "USER_DISABLED"
	CodeUserAuthFailed             = "USER_AUTH_FAILED"
	CodeUserInvalidInfo            = "USER_INVALID_INFO"
	CodeUserQueryFailed            = "USER_QUERY_FAILED"
	CodeUserWxOpenIDEmpty          = "USER_WX_OPENID_EMPTY"
	CodeUserAvatarEmpty            = "USER_AVATAR_EMPTY"
	CodeUserWxTelephoneParseFailed = "USER_WX_TELEPHONE_PARSE_FAILED"
	CodeUserRegisterFailed         = "USER_REGISTER_FAILED"
	CodeUserTelephoneBound         = "USER_TELEPHONE_BOUND"
)

// 预定义用户领域错误
var (
	ErrUserNotFound                 = domainErrors.NewUserError(CodeUserNotFound, "用户不存在", nil)
	ErrUserAlreadyExists            = domainErrors.NewUserError(CodeUserAlreadyExists, "用户已存在", nil)
	ErrUserDisabled                 = domainErrors.NewUserError(CodeUserDisabled, "用户已被禁用", nil)
	ErrUserAuthFailed               = domainErrors.NewUserError(CodeUserAuthFailed, "用户认证失败", nil)
	ErrUserInvalidInfo              = domainErrors.NewUserError(CodeUserInvalidInfo, "用户信息无效", nil)
	ErrUserQueryFailed              = domainErrors.NewUserError(CodeUserQueryFailed, "用户信息查询失败", nil)
	ErrUserWxOpenIDEmpty            = domainErrors.NewUserError(CodeUserWxOpenIDEmpty, "OpenID为空", nil)
	ErrUserAvatarEmpty              = domainErrors.NewUserError(CodeUserAvatarEmpty, "头像地址为空", nil)
	ErrUserWxTelephoneParsingFailed = domainErrors.NewUserError(CodeUserWxTelephoneParseFailed, "微信手机号解析失败", nil)
	ErrUserRegisterFailed           = domainErrors.NewUserError(CodeUserRegisterFailed, "用户注册失败", nil)
	ErrUserTelephoneBound           = domainErrors.NewUserError(CodeUserTelephoneBound, "当前手机号已绑定其它微信用户", nil)
)
