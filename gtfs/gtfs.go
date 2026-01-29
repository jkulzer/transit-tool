package gtfs

import (
	"github.com/jkulzer/transit-tool/db"
	"github.com/jkulzer/transit-tool/env"

	"github.com/jamespfennell/gtfs"

	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/rs/zerolog/log"

	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"time"
)

type StationService struct {
	// IFOPT code
	// https://en.wikipedia.org/wiki/Identification_of_Fixed_Objects_in_Public_Transport
	// could potentially also be something different in different systems, in the VBB GTFS it's IFOPT though
	StopID          string
	ERoutes         map[GtfsRouteID]ExtendedRoute // string is gtfs route id
	RealtimeTripMap map[GtfsTripID]gtfs.Trip
}

type GtfsTripID string

type GtfsRouteID string

type ExtendedRoute struct {
	StopTimesDirectionTrue  []ExtendedStopTime
	StopTimesDirectionFalse []ExtendedStopTime
	StopTimesNoDirection    []ExtendedStopTime
}

type ExtendedStopTime struct {
	StopTime gtfs.ScheduledStopTime
	RTTrip   gtfs.Trip
}

func QueryForDeparture(env *env.Env, stopName string) StationService {
	currentTime := time.Now()
	staticData, err := GetStaticData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	realtimeData, err := getRealtimeData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	log.Debug().Msg("finished parsing data (" + fmt.Sprint(time.Since(currentTime)) + ")")

	service := StationService{}
	service.ERoutes = make(map[GtfsRouteID]ExtendedRoute)

	activeTrips := getActiveTrips(staticData.Trips, currentTime)

	for _, trip := range activeTrips {

		for _, stopTime := range trip.StopTimes {
			// only actually try to match trips which stop at the station we are searching at
			if stopTime.Stop.Name == stopName {

				// the trip field doesn't get populated correctly so i'm doing it myself
				stopTime.Trip = &trip

				extendedRoute := service.ERoutes[GtfsRouteID(stopTime.Trip.Route.Id)]
				extendedStopTime := ExtendedStopTime{StopTime: stopTime}

				// associating realtime trip with scheduled trip
			rtTripLoop:
				for _, rtTrip := range realtimeData.Trips {
					// checks if the IDs of the static and realtime trip are identical
					// the scheduleRelationship thing is if the trip isn't scheduled and only exists in realtime. this is probably a bullshit solution and will have to be implemented more carefully
					// if rtTrip.ID.ID == trip.ID || rtTrip.ID.ScheduleRelationship != 0 {
					if rtTrip.ID.ID == trip.ID {
						log.Debug().Msg("trip id: " + trip.ID + " with route " + trip.Route.ShortName + " matched RT trip with id " + rtTrip.ID.ID + " and has relationship " + fmt.Sprint(rtTrip.ID.ScheduleRelationship))
						// log.Debug().Msg(fmt.Sprint(tripCurrentlyRunning(&trip, currentTime)))
						extendedStopTime.RTTrip = rtTrip
						// extendedStopTime = ExtendedStopTime{StopTime: stopTime, RTTrip: rtTrip}
						break rtTripLoop
					}
				}

				// matching directions
				switch trip.DirectionId {
				case gtfs.DirectionID_Unspecified:
					extendedRoute.StopTimesNoDirection = append(extendedRoute.StopTimesNoDirection, extendedStopTime)
				case gtfs.DirectionID_True:
					extendedRoute.StopTimesDirectionTrue = append(extendedRoute.StopTimesDirectionTrue, extendedStopTime)
				case gtfs.DirectionID_False:
					extendedRoute.StopTimesDirectionFalse = append(extendedRoute.StopTimesDirectionFalse, extendedStopTime)
				}
				service.ERoutes[GtfsRouteID(stopTime.Trip.Route.Id)] = extendedRoute
			}
		}
	}
	log.Debug().Msg("finished creating list(" + fmt.Sprint(time.Since(currentTime)) + ")")
	for key, eRoute := range service.ERoutes {
		eRoute.StopTimesDirectionTrue = sortExtendedStopTimes(eRoute.StopTimesDirectionTrue)
		eRoute.StopTimesDirectionFalse = sortExtendedStopTimes(eRoute.StopTimesDirectionFalse)
		eRoute.StopTimesNoDirection = sortExtendedStopTimes(eRoute.StopTimesNoDirection)
		service.ERoutes[key] = eRoute
	}
	log.Debug().Msg("finished departure query (" + fmt.Sprint(time.Since(currentTime)) + ")")
	return service
}

func GetStaticData(env *env.Env) (*gtfs.Static, error) {
	if env.GtfsStaticData != nil {
		return env.GtfsStaticData, nil
	}
	gtfsSource, err := db.GetGtfsDatasource(env)
	if err != nil {
		log.Err(err).Msg("failed getting gtfs datasource")
		return &gtfs.Static{}, err
	}
	staticGtfsPath := env.App.Storage().RootURI().Path() + "/staticGtfs.zip"
	if _, err := os.Stat(staticGtfsPath); errors.Is(err, os.ErrNotExist) {
		log.Trace().Msg("static gtfs data not cached")
		downloadedFile, err := os.Create(staticGtfsPath)
		if err != nil {
			return &gtfs.Static{}, err
		}
		defer downloadedFile.Close()

		resp, err := http.Get(gtfsSource.StaticUrl)
		if err != nil {
			return &gtfs.Static{}, err
		}
		if resp.StatusCode != http.StatusOK {
			return &gtfs.Static{}, fmt.Errorf("bad status: %s", resp.Status)
		}

		// Writer the body to file
		_, err = io.Copy(downloadedFile, resp.Body)
		if err != nil {
			return &gtfs.Static{}, err
		}
	} else {
		log.Trace().Msg("getting static gtfs data from cache")
	}
	file, err := os.Open(staticGtfsPath)
	if err != nil {
		return &gtfs.Static{}, err
	}
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return &gtfs.Static{}, err
	}
	staticData, err := gtfs.ParseStatic(fileContent, gtfs.ParseStaticOptions{})
	if err != nil {
		return &gtfs.Static{}, err
	}
	env.GtfsStaticData = staticData
	optimizeStaticData(env, staticData)
	return staticData, nil
}

func getRealtimeData(env *env.Env) (*gtfs.Realtime, error) {
	if env.GtfsRealtimeData != nil {
		return env.GtfsRealtimeData, nil
	}
	gtfsSource, err := db.GetGtfsDatasource(env)
	if err != nil {
		log.Err(err).Msg("failed getting gtfs datasource")
		return &gtfs.Realtime{}, err
	}
	realtimeGtfsPath := env.App.Storage().RootURI().Path() + "/realtimeGtfs.bin"
	if _, err := os.Stat(realtimeGtfsPath); errors.Is(err, os.ErrNotExist) {
		// if true {
		log.Trace().Msg("realtime gtfs data not cached")
		downloadedFile, err := os.Create(realtimeGtfsPath)
		if err != nil {
			return &gtfs.Realtime{}, err
		}
		defer downloadedFile.Close()

		resp, err := http.Get(gtfsSource.RealtimeUrl)
		if err != nil {
			return &gtfs.Realtime{}, err
		}
		if resp.StatusCode != http.StatusOK {
			return &gtfs.Realtime{}, fmt.Errorf("bad status: %s", resp.Status)
		}

		// Writer the body to file
		_, err = io.Copy(downloadedFile, resp.Body)
		if err != nil {
			return &gtfs.Realtime{}, err
		}
	} else {
		log.Trace().Msg("getting static gtfs data from cache")
	}
	file, err := os.Open(realtimeGtfsPath)
	if err != nil {
		return &gtfs.Realtime{}, err
	}
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return &gtfs.Realtime{}, err
	}
	realtimeData, err := gtfs.ParseRealtime(fileContent, &gtfs.ParseRealtimeOptions{})
	if err != nil {
		return &gtfs.Realtime{}, err
	}
	env.GtfsRealtimeData = realtimeData
	return realtimeData, nil
}

func optimizeStaticData(envVar *env.Env, staticData *gtfs.Static) {
	envVar.GtfsStaticOptimized = env.GtfsStaticOptimized{}
	envVar.GtfsStaticOptimized.StopTimesByStop = make(map[string][]*gtfs.ScheduledStopTime)
	for _, stop := range staticData.Stops {
		envVar.GtfsStaticOptimized.Stops = append(envVar.GtfsStaticOptimized.Stops, &stop)
	}

	for _, scheduledTrip := range staticData.Trips {
		for _, scheduledStopTime := range scheduledTrip.StopTimes {
			stopTimeStopId := scheduledStopTime.Stop.Id
			envVar.GtfsStaticOptimized.StopTimesByStop[stopTimeStopId] = append(envVar.GtfsStaticOptimized.StopTimesByStop[stopTimeStopId], &scheduledStopTime)
		}
	}
}

func tripPastMidnight(trip *gtfs.ScheduledTrip) bool {
	if trip.StopTimes[len(trip.StopTimes)-1].ArrivalTime > time.Hour*24 {
		return true
	} else {
		return false
	}
}

// could possibly be optimized further with early returns
func tripCurrentlyRunning(trip *gtfs.ScheduledTrip, currentTime time.Time) bool {

	var isActiveService bool

	switch currentTime.Weekday() {
	case time.Monday:
		if trip.Service.Monday || (tripPastMidnight(trip) && trip.Service.Sunday) {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Tuesday:
		if trip.Service.Tuesday || (tripPastMidnight(trip) && trip.Service.Monday) {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Wednesday:
		if trip.Service.Wednesday || (tripPastMidnight(trip) && trip.Service.Tuesday) {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Thursday:
		if trip.Service.Thursday || (tripPastMidnight(trip) && trip.Service.Wednesday) {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Friday:
		if trip.Service.Friday || (tripPastMidnight(trip) && trip.Service.Thursday) {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Saturday:
		if trip.Service.Saturday || (tripPastMidnight(trip) && trip.Service.Friday) {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Sunday:
		if trip.Service.Sunday || (tripPastMidnight(trip) && trip.Service.Saturday) {
			isActiveService = true
		} else {
			isActiveService = false
		}
	}

	for _, addedDate := range trip.Service.AddedDates {
		if addedDate.Year() == currentTime.Year() && addedDate.Month() == currentTime.Month() && addedDate.Day() == currentTime.Day() {
			isActiveService = true
			return true
		}
	}
	for _, removedDate := range trip.Service.RemovedDates {
		if removedDate.Year() == currentTime.Year() && removedDate.Month() == currentTime.Month() && removedDate.Day() == currentTime.Day() {
			isActiveService = false
			return false
		}
	}
	return isActiveService
}

func sortExtendedStopTimes(stopTimes []ExtendedStopTime) []ExtendedStopTime {
	slices.SortFunc(stopTimes, func(a, b ExtendedStopTime) int {
		return int(a.StopTime.DepartureTime - b.StopTime.DepartureTime)
	})
	return stopTimes
}

func SearchStopList(searchTerm string, env *env.Env) fuzzy.Ranks {
	staticData, err := GetStaticData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}

	var stopNameList []string

	for _, stop := range staticData.Stops {
		if stopIsTopLevel(stop) {
			stopNameList = append(stopNameList, stop.Name)
		}
	}

	slices.Sort(stopNameList)
	stopNameList = slices.Compact(stopNameList)

	rankedStopList := fuzzy.RankFindNormalizedFold(searchTerm, stopNameList)
	sort.Sort(rankedStopList)

	return rankedStopList
}

func getActiveTrips(trips []gtfs.ScheduledTrip, currentTime time.Time) []gtfs.ScheduledTrip {
	var activeTrips []gtfs.ScheduledTrip
	for _, trip := range trips {
		if tripCurrentlyRunning(&trip, currentTime) {
			activeTrips = append(activeTrips, trip)
		}
	}
	return activeTrips
}

func GtfsDurationToTime(durationValue time.Duration) time.Time {
	currentTime := time.Now()
	timeValue := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(durationValue)
	return timeValue
}
