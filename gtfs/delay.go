package gtfs

import (
	"fmt"
	"time"

	"github.com/jamespfennell/gtfs"

	"github.com/rs/zerolog/log"
)

func ProcessStopTimeUpdate(stopTimeUpdate gtfs.StopTimeUpdate, scheduledStopTime gtfs.ScheduledStopTime, currentTime time.Time) (bool, time.Duration) {

	departureDuration := scheduledStopTime.DepartureTime
	departureTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(departureDuration)

	var isDelayed bool
	var delay time.Duration
	if stopTimeUpdate.Departure.Delay != nil {
		if *stopTimeUpdate.Departure.Delay > 0 {
			isDelayed = true
		} else {
			isDelayed = false
		}
		delay = *stopTimeUpdate.Departure.Delay

		// if the delay is stored as a time of departure/arrival
	} else if stopTimeUpdate.Departure.Time != nil {
		// fmt.Println(scheduledStopTime.Trip.Route)
		// route := *scheduledStopTime.Trip.Route
		agencyTimezone := time.Now().Location()
		log.Debug().Msg("agency location is " + fmt.Sprint(agencyTimezone))
		departureDate := stopTimeUpdate.Departure.Time
		timezoneCorrectedDate := time.Date(departureDate.Year(), departureDate.Month(), departureTime.Day(), departureTime.Hour(), departureTime.Minute(), departureTime.Second(), departureTime.Nanosecond(), agencyTimezone)
		differenceToScheduled := departureTime.Sub(timezoneCorrectedDate)
		if differenceToScheduled > 0 {
			isDelayed = true
		} else {
			isDelayed = false
		}
		delay = differenceToScheduled
	}
	return isDelayed, delay
}
