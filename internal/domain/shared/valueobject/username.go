package valueobject

import (
	"regexp"

	"github.com/dysodeng/app/internal/domain/shared/errors"
)

// Username 用户名值对象
type Username struct {
	value string
}

func NewUsername(value string) (Username, error) {
	u := Username{
		value: value,
	}
	if err := u.Validate(); err != nil {
		return Username{}, err
	}
	return u, nil
}

func (u Username) Validate() error {
	if u.value == "" {
		return errors.ErrSharedUsernameEmpty
	}
	if len(u.value) < 5 || len(u.value) > 50 {
		return errors.ErrSharedUsernameLengthMismatch
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, u.value)
	if !matched {
		return errors.ErrSharedUsernameInvalid
	}
	return nil
}

// Value 获取用户名值
func (u Username) Value() string {
	return u.value
}
