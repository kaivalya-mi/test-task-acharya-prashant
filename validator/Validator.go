package validator

import (
	"regexp"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	u "test-task/shared/common"
)

type IAPIValidatorService interface {
	ValidateStruct(req interface{}, name string) (string, bool)
}
type APIValidator struct{}

func NewAPIValidatorService() IAPIValidatorService {
	return &APIValidator{}
}

func (uv *APIValidator) ValidateStruct(req interface{}, key string) (string, bool) {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		valErrs := err.(validator.ValidationErrors)
		for _, v := range valErrs {
			fieldName := strings.Replace(strings.Replace(v.Namespace(), key+".", "", 1), ".", " ", 3)
			reg, _ := regexp.Compile("[^A-Z`[]]+")
			fieldName = strings.Replace(reg.ReplaceAllString(fieldName, ""), "[", "", 2)
			errorString := u.GetError(fieldName, v.Tag())
			return errorString, false
		}
	}
	return "", true
}
