package validate

import (
	localeEN "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	en := localeEN.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	validate = validator.New()
	validate.SetTagName("v")

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}
}

func Validated(pointer any) error {
	return validate.Struct(pointer)
}

func Translate(err validator.FieldError) string {
	return err.Translate(trans)
}

type FormRequestValidator struct {
}

func (cv *FormRequestValidator) Validate(pointer any) error {
	return Validated(pointer)
}
