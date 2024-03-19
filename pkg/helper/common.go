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
