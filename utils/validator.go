package utils

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Validator *validator.Validate
	trans ut.Translator
)

//NewValidator 初始化验证器
func InitValidator()  {
	zhCh := zh.New()
	uni := ut.New(zhCh)
	trans, _ = uni.GetTranslator("zh")
	Validator = validator.New()
	zhTranslations.RegisterDefaultTranslations(Validator, trans)
	//Validator.RegisterValidation("is-awesome", ValidateMyVal)
}


//ErrorTranslate 英文错误信息翻译成中文
func ErrorTranslate(err error) interface{} {
	return err.(validator.ValidationErrors).Translate(trans)
}

//func ValidateMyVal(fl validator.FieldLevel) bool {
//	fmt.Println()
//	return fl.Field().String() == "awesome"
//}