package errors

import domainErrors "github.com/dysodeng/app/internal/domain/shared/errors"

const (
	CodePassportUserGrantTypeInvalid      = "PASSPORT_USER_GRANT_TYPE_INVALID"
	CodeLoginFailed                       = "PASSPORT_LOGIN_FAILED"
	CodeLoginUserTypeInvalid              = "PASSPORT_LOGIN_USER_TYPE_INVALID"
	CodePassportGetWxUserFailed           = "PASSPORT_GET_WX_USER_FAILED"
	CodeTokenInvalid                      = "PASSPORT_TOKEN_INVALID"
	CodeBizTokenCannotUsedForRefreshToken = "PASSPORT_BIZ_TOKEN_CANNOT_USE_FOR_REFRESH_TOKEN"
	CodeRefreshTokenCannotUsedForBizToken = "PASSPORT_REFRESH_TOKEN_CANNOT_USE_FOR_BIZ_TOKEN"
	CodeAdminUsernameQueryFailed          = "PASSPORT_ADMIN_USER_QUERY_FAILED"
	CodeAdminUsernameInvalid              = "PASSPORT_ADMIN_USER_INVALID"
	CodeAdminPasswordInvalid              = "PASSPORT_ADMIN_PASSWORD_INVALID"
)

var (
	ErrPassportUserGrantTypeInvalid      = domainErrors.NewPassportError(CodePassportUserGrantTypeInvalid, "登录方式错误", nil)
	ErrLoginFailed                       = domainErrors.NewPassportError(CodeLoginFailed, "登录失败", nil)
	ErrLoginUserTypeInvalid              = domainErrors.NewPassportError(CodeLoginUserTypeInvalid, "登录用户类型错误", nil)
	ErrPassportGetWxUserFailed           = domainErrors.NewPassportError(CodePassportGetWxUserFailed, "微信用户信息获取失败", nil)
	ErrTokenInvalid                      = domainErrors.NewPassportError(CodeTokenInvalid, "Token无效", nil)
	ErrBizTokenCannotUsedForRefreshToken = domainErrors.NewPassportError(CodeBizTokenCannotUsedForRefreshToken, "业务token不能用于刷新token", nil)
	ErrRefreshTokenCannotUsedForBizToken = domainErrors.NewPassportError(CodeRefreshTokenCannotUsedForBizToken, "刷新token不能用于业务请求", nil)
	ErrAdminUsernameQueryFailed          = domainErrors.NewPassportError(CodeAdminUsernameQueryFailed, "管理员信息查询失败", nil)
	ErrAdminUsernameNotFound             = domainErrors.NewPassportError(CodeAdminUsernameInvalid, "登录账号不正确", nil)
	ErrAdminPasswordInvalid              = domainErrors.NewPassportError(CodeAdminPasswordInvalid, "登录密码错误", nil)
)
