package utils

import (
	"fmt"
	"time"
)

// calculateTimeDiff takes two times in "15:04" format and returns the difference as "HH:MM"
func CalculateTimeDiff(start, end string) (string, error) {
	// Parse the input times using a fixed date (because time.Parse needs a full time)
	layout := "15:04"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		return "", fmt.Errorf("invalid start time: %v", err)
	}

	endTime, err := time.Parse(layout, end)
	if err != nil {
		return "", fmt.Errorf("invalid end time: %v", err)
	}

	// Calculate the duration
	duration := endTime.Sub(startTime)

	// Handle negative durations (e.g., if end < start, assume next day)
	if duration < 0 {
		duration += 24 * time.Hour
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%02d:%02d", hours, minutes), nil
}
