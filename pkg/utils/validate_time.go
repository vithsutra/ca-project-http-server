package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateTime(fl validator.FieldLevel) bool {
	inputTime := fl.Field().String()
	_, err := time.Parse("15:04", inputTime)
	return err == nil
}
