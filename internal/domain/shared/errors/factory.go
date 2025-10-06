package errors

// 领域名称常量
const (
	DomainCommon   = "common"
	DomainUser     = "user"
	DomainFile     = "file"
	DomainPassport = "passport"
)

// NewCommonError 创建通用领域错误
func NewCommonError(code, message string, err error) *DomainError {
	return NewDomainError(DomainCommon, code, message, err)
}

// NewUserError 创建用户领域错误
func NewUserError(code, message string, err error) *DomainError {
	return NewDomainError(DomainUser, code, message, err)
}

// NewFileError 创建文件领域错误
func NewFileError(code, message string, err error) *DomainError {
	return NewDomainError(DomainFile, code, message, err)
}

// NewPassportError 创建通行证领域错误
func NewPassportError(code, message string, err error) *DomainError {
	return NewDomainError(DomainPassport, code, message, err)
}
