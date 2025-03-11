package utils

import (
	"errors"
	"time"
)

func CompareDates(date1 string, date2 string) error {
	date1Parsed, err := time.Parse("2006-01-02", date1)

	if err != nil {
		return errors.New("invalid date formates")
	}

	date2Parsed, err := time.Parse("2006-01-02", date2)

	if err != nil {
		return errors.New("invalid date formates")
	}

	if date2Parsed.After(date1Parsed) {
		return nil
	} else if date1Parsed.Equal(date2Parsed) {
		return nil
	} else {
		return errors.New("dates order mismatch")
	}

}
