package valueobject

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/shared/validator"
)

// Email 邮箱地址值对象
type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	vo := Email{
		email: email,
	}
	if err := vo.Validate(); err != nil {
		return Email{}, err
	}
	return vo, nil
}

func (e Email) String() string {
	return e.email
}

func (e Email) Validate() error {
	if e.email == "" {
		return errors.New("邮箱地址为空")
	}
	if !validator.IsEmail(e.email) {
		return errors.New("邮箱地址格式错误")
	}
	return nil
}
