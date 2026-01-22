package gtfs

import (
	"github.com/jkulzer/transit-tool/env"

	"github.com/jamespfennell/gtfs"

	"github.com/rs/zerolog/log"

	"strings"
	"time"
)

type Journey struct {
	Length          time.Duration
	MemberStopTimes []gtfs.ScheduledStopTime
	MemberStops     []gtfs.Stop
	MemberTrips     []gtfs.ScheduledTrip
}

func CalculateJourney(env *env.Env, departureTime time.Time, departureStation, arrivalStation string) Journey {
	journeys := make(map[string]Journey)
	stationService := QueryForDeparture(env, departureStation)

	journeys[departureStation] = Journey{Length: 0}

	for _, eRoute := range stationService.ERoutes {
		if len(eRoute.StopTimesDirectionTrue) > 0 {
			processStopTimes(eRoute.StopTimesDirectionTrue, departureTime, &journeys)
		}
		if len(eRoute.StopTimesDirectionFalse) > 0 {
			processStopTimes(eRoute.StopTimesDirectionFalse, departureTime, &journeys)
		}
		if len(eRoute.StopTimesNoDirection) > 0 {
			processStopTimes(eRoute.StopTimesNoDirection, departureTime, &journeys)
		}
	}

	log.Debug().Msg("arrival station: :" + arrivalStation)

	var arrivalStopID string

	for _, stop := range env.GtfsStaticData.Stops {
		if stop.Name == arrivalStation {
			arrivalStopID = stop.Root().Id
			break
		}
	}

	log.Debug().Msg("arrival stop id: " + arrivalStopID)

	for _, journey := range journeys {
		if len(journey.MemberStops) > 0 {
			if journey.MemberStops[len(journey.MemberStops)-1].Root().Id == arrivalStopID {
				return journey
			}
		}
	}
	return Journey{}
}

// TODO: this is full of copy-pasting from widgets/departureChip.go, fix it
// TODO: this function sucks, document it
func processStopTimes(
	extendedStopTimes []ExtendedStopTime,
	requestedDepartureTime time.Time,
	journeys *map[string]Journey, // map of journeys, key is stop id
) {

	for _, extendedStopTime := range extendedStopTimes {
		// time train leaves from stop
		departureTime := GtfsDurationToTime(extendedStopTime.StopTime.DepartureTime)

	stopTimeUpdateLoop:
		for _, stopTimeUpdate := range extendedStopTime.RTTrip.StopTimeUpdates {
			// matches stop time with realtime trip stop time
			if strings.Contains(extendedStopTime.StopTime.Stop.Root().Id, *stopTimeUpdate.StopID) {
				// gets the current delay from the realtime stop time
				_, _, departureDelay, _ := ProcessStopTimeUpdate(stopTimeUpdate, extendedStopTime.StopTime, departureTime)
				departureTime.Add(departureDelay)
				// found matching rt stop time so breaking loop
				break stopTimeUpdateLoop
			}
		}

		// gets the current journey to the stop in the list of stop times
		currentlyBestJourney := (*journeys)[extendedStopTime.StopTime.Stop.Root().Id]
		if currentlyBestJourney.Length != time.Duration(0) {
			// 	if currentlyBestJourney.Length >  {
			//
			// }
		} else {
		}

		for _, stopTime := range extendedStopTime.StopTime.Trip.StopTimes {
			arrivalTime := GtfsDurationToTime(stopTime.ArrivalTime)

			if arrivalTime.After(departureTime) {
				journeyForStop := (*journeys)[stopTime.Stop.Root().Id]
				journeyForStop.Length = arrivalTime.Sub(departureTime)
				// journeyForStop.MemberStops = append(journeyForStop.MemberStops, *stopTime.Stop)
				journeyForStop.MemberStopTimes = []gtfs.ScheduledStopTime{stopTime}
				journeyForStop.MemberStops = []gtfs.Stop{*stopTime.Stop}
				journeyForStop.MemberTrips = []gtfs.ScheduledTrip{*extendedStopTime.StopTime.Trip}
				(*journeys)[stopTime.Stop.Root().Id] = journeyForStop
			}
		}
	}
}
