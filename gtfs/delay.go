package gtfs

import (
	"time"

	"image/color"

	"github.com/jkulzer/transit-tool/colors"

	"github.com/jamespfennell/gtfs"
)

type DelayStatus int

const (
	OnTime DelayStatus = iota
	Delayed
	Early
)

func ProcessStopTimeUpdate(stopTimeUpdate gtfs.StopTimeUpdate, scheduledStopTime gtfs.ScheduledStopTime, currentTime time.Time) (arrivalDelay time.Duration, arrivalDelayStatus DelayStatus, departureDelay time.Duration, departureDelayStatus DelayStatus) {

	if stopTimeUpdate.Arrival.Delay != nil {
		departureDelay, departureDelayStatus = processDelay(stopTimeUpdate.Arrival)
	} else if stopTimeUpdate.Arrival.Time != nil {
		departureDelay, departureDelayStatus = processStopEventTime(stopTimeUpdate.Arrival, currentTime, scheduledStopTime.ArrivalTime)
	}
	if stopTimeUpdate.Departure.Delay != nil {
		departureDelay, departureDelayStatus = processDelay(stopTimeUpdate.Departure)
	} else if stopTimeUpdate.Departure.Time != nil {
		departureDelay, departureDelayStatus = processStopEventTime(stopTimeUpdate.Departure, currentTime, scheduledStopTime.DepartureTime)
	}

	return arrivalDelay, arrivalDelayStatus, departureDelay, departureDelayStatus
}

// if the delay is stored as a time of departure/arrival
func processDelay(stopTimeEvent *gtfs.StopTimeEvent) (time.Duration, DelayStatus) {
	var delayStatus DelayStatus
	var delay time.Duration
	if *stopTimeEvent.Delay > 0 {
		delayStatus = Delayed
	} else if *stopTimeEvent.Delay < 0 {
		delayStatus = Early
	} else {
		delayStatus = OnTime
	}
	delay = *stopTimeEvent.Delay

	return delay, delayStatus
}

// if the delay is stored as a time of departure/arrival
func processStopEventTime(stopTimeEvent *gtfs.StopTimeEvent, currentTime time.Time, staticEventDuration time.Duration) (time.Duration, DelayStatus) {
	staticEventTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(staticEventDuration)

	var delayStatus DelayStatus
	var delay time.Duration

	// TODO: replace it with actually reading the agency timezone, this assumes the user is in the same timezone as the agency
	agencyTimezone := time.Now().Location()
	eventDate := stopTimeEvent.Time
	timezoneCorrectedDate := time.Date(eventDate.Year(), eventDate.Month(), staticEventTime.Day(), staticEventTime.Hour(), staticEventTime.Minute(), staticEventTime.Second(), staticEventTime.Nanosecond(), agencyTimezone)
	differenceToScheduled := staticEventTime.Sub(timezoneCorrectedDate)
	if differenceToScheduled > 0 {
		delayStatus = Delayed
	} else if differenceToScheduled < 0 {
		delayStatus = Early
	} else {
		delayStatus = OnTime
	}
	delay = differenceToScheduled

	return delay, delayStatus
}

func GetDelayColor(delay time.Duration, delayStatus DelayStatus, hasRealtimeData bool) (string, color.Color) {
	var delayString string
	var delayColor color.Color
	if hasRealtimeData {
		delayString = delay.String()
		switch delayStatus {
		case Delayed:
			delayColor = colors.Red()
			delayString = "+" + delayString
		case Early:
			delayColor = colors.Blue()
		case OnTime:
			delayColor = colors.Green()
		}
	} else {
		delayColor = color.White
		delayString = ""
	}

	return delayString, delayColor
}
