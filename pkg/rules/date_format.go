package rules

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var dateRegexp = regexp.MustCompile(`^(0[1-9]|1[0-2])-\d{4}$`)

func DateFormat(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if value == "" {
		return true
	}
	return dateRegexp.MatchString(value)
}
