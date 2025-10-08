package valueobject

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/shared/validator"
)

// Telephone 手机号值对象
type Telephone struct {
	value string
}

func NewTelephone(telephone string) (Telephone, error) {
	telephone = strings.TrimSpace(telephone)
	t := Telephone{value: telephone}
	if err := t.Validate(); err != nil {
		return Telephone{}, err
	}
	return t, nil
}

func NewEmptyTelephone() Telephone {
	return Telephone{}
}

func (t Telephone) String() string {
	return t.Value()
}

func (t Telephone) Value() string {
	return t.value
}

func (t Telephone) Validate() error {
	if t.value == "" {
		return errors.New("手机号为空")
	}
	if !validator.IsMobile(t.value) {
		return errors.New("手机号格式错误")
	}
	return nil
}
