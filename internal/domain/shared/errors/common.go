package errors

// 通用错误码
const (
	// PrefixCommon 通用错误码前缀
	PrefixCommon = "COMMON"

	// CodeCommonUnknownError 未知错误
	CodeCommonUnknownError = "COMMON_UNKNOWN_ERROR"
	// CodeCommonValidationError 参数验证错误
	CodeCommonValidationError = "COMMON_VALIDATION_ERROR"
	// CodeCommonNotFound 资源不存在
	CodeCommonNotFound = "COMMON_NOT_FOUND"
	// CodeCommonAlreadyExists 资源已存在
	CodeCommonAlreadyExists = "COMMON_ALREADY_EXISTS"
	// CodeCommonOperationFailed 操作失败
	CodeCommonOperationFailed = "COMMON_OPERATION_FAILED"
	// CodeCommonPermissionDenied 权限不足
	CodeCommonPermissionDenied = "COMMON_PERMISSION_DENIED"
	// CodeCommonUnauthorized 未授权
	CodeCommonUnauthorized = "COMMON_UNAUTHORIZED"
	// CodeCommonInternalError 系统内部错误
	CodeCommonInternalError = "COMMON_INTERNAL_ERROR"
)

// 预定义通用错误
var (
	ErrCommonUnknown          = NewCommonError(CodeCommonUnknownError, "未知错误", nil)
	ErrCommonValidation       = NewCommonError(CodeCommonValidationError, "参数验证错误", nil)
	ErrCommonNotFound         = NewCommonError(CodeCommonNotFound, "资源不存在", nil)
	ErrCommonAlreadyExists    = NewCommonError(CodeCommonAlreadyExists, "资源已存在", nil)
	ErrCommonOperationFailed  = NewCommonError(CodeCommonOperationFailed, "操作失败", nil)
	ErrCommonPermissionDenied = NewCommonError(CodeCommonPermissionDenied, "权限不足", nil)
	ErrCommonUnauthorized     = NewCommonError(CodeCommonUnauthorized, "未授权", nil)
	ErrCommonInternalError    = NewCommonError(CodeCommonInternalError, "系统内部错误", nil)
)
