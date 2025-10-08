package errors

// 用户名错误码
const (
	CodeSharedUsernameEmpty          = "SHARED_USERNAME_EMPTY"
	CodeSharedUsernameLengthMismatch = "SHARED_USERNAME_LENGTH_MISMATCH"
	CodeSharedUsernameInvalid        = "SHARED_USERNAME_INVALID"
	CodeSharedUsernameAlreadyExists  = "SHARED_USERNAME_ALREADY_EXISTS"
)

// 密码错误码
const (
	CodeSharedPasswordEmpty                  = "SHARED_PASSWORD_EMPTY"
	CodeSharedPasswordLengthMismatch         = "SHARED_PASSWORD_LENGTH_MISMATCH"
	CodeSharedPasswordInvalid                = "SHARED_PASSWORD_INVALID"
	CodeSharedPasswordGenerateFailed         = "SHARED_PASSWORD_GENERATE_FAILED"
	CodeSharedPasswordTooWeak                = "SHARED_PASSWORD_TOO_WEAK"
	CodeSharedPasswordComplexityInsufficient = "SHARED_PASSWORD_COMPLEXITY_INSUFFICIENT"
	CodeSharedPasswordHasConsecutiveChars    = "SHARED_PASSWORD_HAS_CONSECUTIVE_CHARS"
	CodeSharedPasswordHasRepeatingChars      = "SHARED_PASSWORD_HAS_REPEATING_CHARS"
)

// 用户名相关错误
var (
	ErrSharedUsernameEmpty          = NewShardError(CodeSharedUsernameEmpty, "登录账号为空", nil)
	ErrSharedUsernameLengthMismatch = NewShardError(CodeSharedUsernameLengthMismatch, "登录账号长度不符", nil)
	ErrSharedUsernameInvalid        = NewShardError(CodeSharedUsernameInvalid, "登录账号格式错误", nil)
	ErrSharedUsernameAlreadyExists  = NewShardError(CodeSharedUsernameAlreadyExists, "登录账号已存在", nil)
)

// 密码相关错误
var (
	ErrSharedPasswordEmpty                  = NewShardError(CodeSharedPasswordEmpty, "登录密码为空", nil)
	ErrSharedPasswordLengthMismatch         = NewShardError(CodeSharedPasswordLengthMismatch, "登录密码长度应在8到128之间", nil)
	ErrSharedPasswordInvalid                = NewShardError(CodeSharedPasswordInvalid, "无效的登录密码", nil)
	ErrSharedPasswordGenerateFailed         = NewShardError(CodeSharedPasswordGenerateFailed, "登录密码生成失败", nil)
	ErrSharedPasswordTooWeak                = NewShardError(CodeSharedPasswordTooWeak, "密码过于简单，请使用更复杂的密码", nil)
	ErrSharedPasswordComplexityInsufficient = NewShardError(CodeSharedPasswordComplexityInsufficient, "密码复杂度不足，至少包含大写字母、小写字母、数字、特殊字符中的3种", nil)
	ErrSharedPasswordHasConsecutiveChars    = NewShardError(CodeSharedPasswordHasConsecutiveChars, "密码不能包含连续字符（如123、abc）", nil)
	ErrSharedPasswordHasRepeatingChars      = NewShardError(CodeSharedPasswordHasRepeatingChars, "密码不能包含过多重复字符", nil)
)
