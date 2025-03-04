package utils

import (
	"fmt"
	"time"
)

func StringToTime(s string, errorMessage string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now().UTC(), fmt.Errorf("%s: %w", errorMessage, err)
	}
	return t, nil
}
