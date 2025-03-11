package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateLatitude(fl validator.FieldLevel) bool {
	location := fl.Field().String()
	re := regexp.MustCompile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`)
	return re.MatchString(location)
}

func ValidateLongitude(fl validator.FieldLevel) bool {
	location := fl.Field().String()
	re := regexp.MustCompile(`^[-+]?((1[0-7]\d|[1-9]?\d)(\.\d+)?|180(\.0+)?)$`)
	return re.MatchString(location)
}
