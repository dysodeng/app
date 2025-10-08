package valueobject

import (
	"strings"
	"unicode"

	"github.com/dysodeng/app/internal/domain/shared/errors"
	"github.com/dysodeng/app/internal/infrastructure/shared/helper"
)

const (
	// MinPasswordLength 密码长度限制
	MinPasswordLength = 8
	MaxPasswordLength = 128

	// MinRequiredCharTypes 密码强度要求
	MinRequiredCharTypes = 3 // 至少包含3种字符类型
)

// 常见弱密码列表
var commonWeakPasswords = []string{
	"password", "123456", "12345678", "qwerty", "abc123", "password123",
	"admin", "root", "user", "guest", "test", "demo", "welcome",
	"123456789", "1234567890", "qwerty123", "password1", "admin123",
}

// Password 登录密码值对象
type Password struct {
	value string // 保存的hash后的密文
}

func NewPassword(value string) (Password, error) {
	if err := validatePasswordStrength(value); err != nil {
		return Password{}, err
	}

	hashPassword, err := helper.GeneratePassword(value)
	if err != nil {
		return Password{}, errors.ErrSharedPasswordGenerateFailed.Wrap(err)
	}

	return Password{
		value: hashPassword,
	}, nil
}

// validatePasswordStrength 验证密码强度
func validatePasswordStrength(password string) error {
	// 检查密码是否为空
	if password == "" {
		return errors.ErrSharedPasswordEmpty
	}

	// 检查密码长度
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return errors.ErrSharedPasswordLengthMismatch
	}

	// 检查是否为常见弱密码
	lowerPassword := strings.ToLower(password)
	for _, weak := range commonWeakPasswords {
		if lowerPassword == weak {
			return errors.ErrSharedPasswordTooWeak
		}
	}

	// 检查字符类型多样性
	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 计算字符类型数量
	charTypeCount := 0
	if hasLower {
		charTypeCount++
	}
	if hasUpper {
		charTypeCount++
	}
	if hasDigit {
		charTypeCount++
	}
	if hasSpecial {
		charTypeCount++
	}

	// 要求至少包含3种字符类型
	if charTypeCount < MinRequiredCharTypes {
		return errors.ErrSharedPasswordComplexityInsufficient
	}

	// 检查连续字符（如123、abc）
	if hasConsecutiveChars(password) {
		return errors.ErrSharedPasswordHasConsecutiveChars
	}

	// 检查重复字符（如aaa、111）
	if hasRepeatingChars(password, 3) {
		return errors.ErrSharedPasswordHasRepeatingChars
	}

	return nil
}

// hasConsecutiveChars 检查是否包含连续字符
func hasConsecutiveChars(password string) bool {
	consecutiveCount := 1
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1]+1 || password[i] == password[i-1]-1 {
			consecutiveCount++
			if consecutiveCount >= 3 {
				return true
			}
		} else {
			consecutiveCount = 1
		}
	}
	return false
}

// hasRepeatingChars 检查是否包含重复字符
func hasRepeatingChars(password string, maxRepeats int) bool {
	charCount := make(map[rune]int)
	for _, char := range password {
		charCount[char]++
		if charCount[char] >= maxRepeats {
			return true
		}
	}
	return false
}

func NewPasswordByHashText(hashPassword string) (Password, error) {
	if hashPassword == "" {
		return Password{}, errors.ErrSharedPasswordEmpty
	}
	return Password{
		value: hashPassword,
	}, nil
}

func NewEmptyPassword() Password {
	return Password{}
}

func (p *Password) Verify(plainPassword string) bool {
	return helper.VerifyPassword(p.value, plainPassword)
}

// Value 获取密码哈希值
func (p *Password) Value() string {
	return p.value
}

// IsEmpty 检查密码是否为空
func (p *Password) IsEmpty() bool {
	return p.value == ""
}

// PasswordStrength 密码强度等级
type PasswordStrength int

const (
	PasswordStrengthWeak PasswordStrength = iota
	PasswordStrengthMedium
	PasswordStrengthStrong
	PasswordStrengthVeryStrong
)

// String 返回密码强度的字符串表示
func (ps PasswordStrength) String() string {
	switch ps {
	case PasswordStrengthWeak:
		return "弱"
	case PasswordStrengthMedium:
		return "中等"
	case PasswordStrengthStrong:
		return "强"
	case PasswordStrengthVeryStrong:
		return "很强"
	default:
		return "未知"
	}
}

// EvaluatePasswordStrength 评估密码强度
func EvaluatePasswordStrength(password string) PasswordStrength {
	if len(password) < 6 {
		return PasswordStrengthWeak
	}

	score := 0

	// 长度评分
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if len(password) >= 16 {
		score++
	}

	// 字符类型评分
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	// 检查是否为常见弱密码
	lowerPassword := strings.ToLower(password)
	for _, weak := range commonWeakPasswords {
		if lowerPassword == weak {
			return PasswordStrengthWeak
		}
	}

	// 根据评分返回强度等级
	switch {
	case score <= 3:
		return PasswordStrengthWeak
	case score <= 5:
		return PasswordStrengthMedium
	case score <= 7:
		return PasswordStrengthStrong
	default:
		return PasswordStrengthVeryStrong
	}
}
