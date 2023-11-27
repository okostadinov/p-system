package validator

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type FormErrors map[string]string

type Validator struct {
	Validate   *validator.Validate
	FormErrors FormErrors
}

var validate *validator.Validate

// sets the validator tag reference to gorilla/schema's tag, add a custom password tag validator, and return the validator
func NewValidator() *Validator {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("schema")
	})

	validate.RegisterValidation("password", passwordValidate)

	return &Validator{Validate: validate, FormErrors: make(FormErrors)}
}

// parses the form and returns whether the validation was successful or not, while storing the errors in the FormErrors map
func (v *Validator) ValidateForm(form interface{}) bool {
	err := v.Validate.Struct(form)

	if err != nil {
		v.FormErrors = make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			v.FormErrors[err.Field()] = v.fetchTagErrorMessage(err.Tag(), err.Param())
		}

		return false
	}

	return true
}

// associates each validation error with a human comprehensible message
func (v *Validator) fetchTagErrorMessage(tag, param string) string {
	switch tag {
	case "required":
		return "required field"
	case "numeric":
		return "invalid format (only numbers allowed)"
	case "len":
		return fmt.Sprintf("invalid amount (requires %v)", param)
	case "alphaunicode":
		return "invalid format (only letters allowed)"
	case "e164":
		return "invalid format (e.g. +359123456789)"
	case "password":
		return "invalid format (requires minimum 8 characters, including letters and numbers)"
	case "eqfield":
		return fmt.Sprintf("field does not equal %s", param)
	case "email":
		return "invalid format (e.g. email@example.com)"
	default:
		return "undefined error"
	}
}

// validates whether a string is at least 8 characters long and includes both letters and numbers
func passwordValidate(fl validator.FieldLevel) bool {
	var (
		hasLetters = false
		hasNumbers = false
	)

	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetters = true
		case unicode.IsNumber(c):
			hasNumbers = true
		}
	}

	return hasLetters && hasNumbers
}
