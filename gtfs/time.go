package gtfs

import (
	"math"
	"time"
)

func DurationToTime(duration time.Duration) time.Time {
	currentTime := time.Now()

	durationHours := int(math.Floor(duration.Hours()))
	durationMinutes := int(math.Floor(math.Mod(duration.Minutes(), 60.0)))
	durationSeconds := int(math.Floor(math.Mod(duration.Seconds(), 60.0)))

	resultingTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), durationHours, durationMinutes, durationSeconds, 0, currentTime.Location())
	return resultingTime
}
