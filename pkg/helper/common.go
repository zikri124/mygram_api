package helper

import (
	"math"
	"time"
)

func CountAge(dob time.Time) uint16 {
	age := time.Since(dob).Hours()
	age = age / 24 / 365
	age = math.Round(age)

	return uint16(age)
}

func ParseStrToTime(timeStr string) (*time.Time, error) {
	time, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		return nil, err
	}

	return &time, nil
}
