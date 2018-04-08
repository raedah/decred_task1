package ctrl

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"strings"
)

func validateStruct(target interface{},ignoreFields ...string) (error,validator.ValidationErrorsTranslations) {
	validate := validator.New()
	enLang := en.New()
	uni := ut.New(enLang)

	trans, _ := uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.StructExcept(target,ignoreFields...)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		tran := errs.Translate(trans)
		return err,convertValidation(tran)
	}
	return nil,nil
}

func convertValidation(validMess validator.ValidationErrorsTranslations) validator.ValidationErrorsTranslations {
	if(validMess == nil){
		return validMess
	}
	newMess := make(validator.ValidationErrorsTranslations)
	for k,mess := range validMess{
		if index := strings.Index(k,".");index != -1{
			k = k[index+1:]
		}
		newMess[k] = mess
	}

	return newMess
}
