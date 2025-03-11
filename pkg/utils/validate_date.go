package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
