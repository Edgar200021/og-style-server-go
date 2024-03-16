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
		str := generateErrMessage(err.Field(), err.Tag(), err.Param())
		errors = append(errors, [2]string{strings.ToLower(err.Field()), str})

	}
	return &errors
}

func generateErrMessage(field, tag, param string) string {
	switch tag {
	case "email":
		return "Invalid email"
	case "gte":
		return fmt.Sprintf("%s must be more than %s symbols", field, param)
	case "required":
		return fmt.Sprintf("%s is required", field)
	default:
		return ""
	}
}
