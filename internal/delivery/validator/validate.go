package validator

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var msgs []string
	for _, ve := range err.(validator.ValidationErrors) {
		msgs = append(msgs, fmt.Sprintf("%s: %s(%s)", ve.Field(), ve.Tag(), ve.Param()))
	}
	return fmt.Errorf("validation failed: %s", strings.Join(msgs, "; "))
}
