package validator

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	trans ut.Translator
)

// InitValidator 初始化验证器
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhTrans := zh.New()
		uni := ut.New(zhTrans, zhTrans)
		trans, _ = uni.GetTranslator("zh")

		// 注册一个函数，获取 struct tag 中的 msg 标签作为错误提示信息
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			msg := fld.Tag.Get("msg")
			if msg == "" {
				msg = fld.Tag.Get("json")
			}
			return msg
		})
	}
}

// TransError 翻译错误信息
func TransError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			return e.Field()
		}
	}
	return err.Error()
}

func Translator() ut.Translator {
	return trans
}
