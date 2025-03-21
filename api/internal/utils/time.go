package utils

import (
	"time"
)

func StringToTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now().UTC(), err
	}
	return t, nil
}
