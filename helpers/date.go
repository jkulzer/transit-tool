package helpers

import (
	"time"
)

func DateIsToday(currentTime, timeToCheck time.Time) bool {
	if currentTime.Year() == timeToCheck.Year() && currentTime.Month() == timeToCheck.Month() && currentTime.Day() == timeToCheck.Day() {
		return true
	} else {
		return false
	}
}

func DateIsYesterday(currentTime, timeToCheck time.Time) bool {
	return DateIsToday(currentTime, timeToCheck.AddDate(0, 0, 1))
}

func DateIsTomorrow(currentTime, timeToCheck time.Time) bool {
	return DateIsToday(currentTime, timeToCheck.AddDate(0, 0, -1))
}
