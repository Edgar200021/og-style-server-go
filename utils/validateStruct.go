package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateStruct(s any) *[][2]string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := [][2]string{}
	for _, err := range err.(validator.ValidationErrors) {
		fmt.Println(err.Field(), err.Tag(), err.Param())
		str := generateErrMessage(err.Field(), err.Tag(), err.Param())
		errors = append(errors, [2]string{strings.ToLower(err.Field()), str})

	}
	return &errors
}

func generateErrMessage(field, tag, param string) string {
	switch tag {
	case "email":
		return "Не корректный эл. адрес"
	case "gte":
		return fmt.Sprintf("должно быть более %s символов", param)
	case "lte":
		return fmt.Sprintf("должно быть меньше %s символов", param)
	case "required":
		return "это поле является обязательным"
	case "oneof":
		return fmt.Sprintf("должно быть один из вариантов %s", strings.Join(strings.Split(param, " "), ","))
	case "min":
		return fmt.Sprintf("Минимальное значение %s", param)
	case "max":
		return fmt.Sprintf("Максимальное значение %s", param)
	case "required_with":
		return fmt.Sprintf("Это поле обязательно для заполнения после выбора значения в поле %s", param)
	case "len":
		return fmt.Sprintf("количество элементов должно быть %s", param)
	default:
		return ""
	}
}
