package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(validate *validator.Validate, data interface{}) map[string]string {
	errors := make(map[string]string)

	if errs := validate.Struct(data); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			msg := fmt.Sprintf("must implement '%s'", err.Tag())
			errors[strings.ToLower(err.Field())] = msg
		}
	}

	return errors
}
