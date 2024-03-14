package types

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type CreateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

func (cu *CreateUser) Validate() *[][2]string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cu)
	if err == nil {
		return nil
	}

	errors := [][2]string{}
	for _, err := range err.(validator.ValidationErrors) {
		str := generateErrMessage(err.Field(), err.Tag(), err.Param())
		errors = append(errors, [2]string{err.Field(), str})

	}
	return &errors
}

type UpdateUser struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=8"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

func (cu *UpdateUser) Validate() *[][2]string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cu)
	if err == nil {
		return nil
	}

	errors := [][2]string{}
	for _, err := range err.(validator.ValidationErrors) {
		str := generateErrMessage(err.Field(), err.Tag(), err.Param())
		errors = append(errors, [2]string{err.Field(), str})

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
