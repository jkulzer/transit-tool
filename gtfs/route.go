package gtfs

import (
	"github.com/jkulzer/transit-tool/env"

	"github.com/jamespfennell/gtfs"

	"github.com/rs/zerolog/log"

	"time"
)

type Journey struct {
	EarliestTimeAtStop time.Duration
	TripToBoardedFrom  gtfs.Stop
}

func CalculateJourney(env *env.Env, departureTime time.Time, departureStation, arrivalStation string, maxTransfers uint) []map[string]Journey {
	journeys := make([]map[string]Journey, maxTransfers+2, maxTransfers+2)
	journeys[0] = make(map[string]Journey)
	journeys[0][departureStation] = Journey{EarliestTimeAtStop: TimeToDuration(departureTime)}

	routeMap := make(map[string]gtfs.Route)
	tripMap := make(map[string]gtfs.ScheduledTrip)
	tripsByRoute := make(map[string][]gtfs.ScheduledTrip)

	for _, route := range env.GtfsStaticData.Routes {
		routeID := route.Id
		routeMap[routeID] = route
	}
	for _, scheduledTrip := range env.GtfsStaticData.Trips {
		tripMap[scheduledTrip.ID] = scheduledTrip
		routeID := scheduledTrip.Route.Id
		tripsByRoute[routeID] = append(tripsByRoute[routeID], scheduledTrip)
	}

	for i := uint(0); i <= maxTransfers; i++ {
		journeys[i+1] = make(map[string]Journey)
		journeys[i+1] = journeys[i]
		for _, tripArray := range tripsByRoute {
		tripLoop:
			for _, scheduledTrip := range tripArray {
				hoppedOn := false
				var stopHoppedOn gtfs.Stop
				for _, scheduledStopTime := range scheduledTrip.StopTimes {
					rootStopID := scheduledStopTime.Stop.Root().Id
					earliestTimeAtStop := journeys[i][rootStopID].EarliestTimeAtStop
					// because gtfs is weird the departure time is the duration from the start of the first day of the trip (yes this duration can be 24h+)
					// TODO: handle 24h+ edge case
					if scheduledStopTime.DepartureTime > earliestTimeAtStop && (earliestTimeAtStop > 1) && hoppedOn == false {
						hoppedOn = true
						stopHoppedOn = *scheduledStopTime.Stop.Root()
						log.Debug().Msg("hopped on trip " + scheduledTrip.Route.ShortName + " to " + scheduledTrip.Headsign + " at " + scheduledStopTime.DepartureTime.String() + " at stop " + scheduledStopTime.Stop.Name)
					}
					if hoppedOn {
						if journeys[i+1][rootStopID].EarliestTimeAtStop > scheduledStopTime.ArrivalTime || journeys[i+1][rootStopID].EarliestTimeAtStop == 0 {
							journeys[i+1][rootStopID] = Journey{
								EarliestTimeAtStop: scheduledStopTime.ArrivalTime,
								TripToBoardedFrom:  stopHoppedOn,
							}
						}
					}
				}
				if hoppedOn {
					break tripLoop
				}
			}
		}
	}

	log.Debug().Msg("arrival station: :" + arrivalStation)

	return journeys
}
