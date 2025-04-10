package utils

import (
	"fmt"
	"strings"
	"time"
)

func SumTimes(times []string) (string, error) {
	totalMinutes := 0

	for _, t := range times {
		parts := strings.Split(t, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid time format: %s", t)
		}

		parsedTime, err := time.Parse("15:04", t)
		if err != nil {
			return "", fmt.Errorf("failed to parse time %s: %v", t, err)
		}

		// Extract hour and minute from parsed time
		hour := parsedTime.Hour()
		minute := parsedTime.Minute()

		totalMinutes += hour*60 + minute
	}

	// Convert total minutes to HH:MM
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	return fmt.Sprintf("%02d:%02d", hours, minutes), nil
}
