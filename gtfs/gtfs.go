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
	"strings"
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

func FindStop(env *env.Env, searchString string) ([]gtfs.Stop, error) {
	staticData, err := getStaticData(env)
	if err != nil {
		return nil, err
	}
	var stopList []gtfs.Stop
	for _, stop := range staticData.Stops {
		if strings.Contains(stop.Name, searchString) {
			if stopIsTopLevel(stop) {
				stopList = append(stopList, stop)
			}
		}
	}
	return stopList, nil
}

// func mapScheduledAndRealtimeTrips(realtimeTrips []gtfs.Trip, scheduledTrips []gtfs.ScheduledTrip) map[GtfsTripID]gtfs.Trip {
// 	tripMap := make(map[GtfsTripID]gtfs.Trip)
// 	for _, scheduledTrip := range scheduledTrips {
// 		for _, realtimeTrip := range realtimeTrips {
// 			if scheduledTrip.ID == realtimeTrip.ID.ID {
// 				tripMap[GtfsTripID(scheduledTrip.ID)] = realtimeTrip
// 				log.Debug().Msg("found a rt and scheduled trip match")
// 				fmt.Println("")
// 				fmt.Println(scheduledTrip.ID)
// 				fmt.Println(realtimeTrip.ID)
// 			}
// 		}
// 	}
// 	return tripMap
// }

func QueryForDeparture(env *env.Env, stopName string) StationService {
	currentTime := time.Now()
	staticData, err := getStaticData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	realtimeData, err := getRealtimeData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	log.Debug().Msg("doing stuff with realtime data")
	// rtScheduledTripMap := mapScheduledAndRealtimeTrips(realtimeData.Trips, staticData.Trips)
	log.Debug().Msg("finished parsing data (" + fmt.Sprint(time.Since(currentTime)) + ")")
	// foundStops, _ := FindStop(env, stopName)
	log.Debug().Msg("finished finding stops (" + fmt.Sprint(time.Since(currentTime)) + ")")

	service := StationService{}
	service.ERoutes = make(map[GtfsRouteID]ExtendedRoute)

	for _, trip := range staticData.Trips {
		for _, rtTrip := range realtimeData.Trips {
			if rtTrip.ID.ID == trip.ID || rtTrip.ID.ScheduleRelationship != 0 {
				if trip.Route.ShortName == "U5" {
					log.Debug().Msg("trip id: " + trip.ID + " with route " + trip.Route.ShortName + " matched RT trip with id " + rtTrip.ID.ID + " and has relationship " + fmt.Sprint(rtTrip.ID.ScheduleRelationship))
					log.Debug().Msg(fmt.Sprint(serviceCurrentlyRunning(trip.Service, currentTime)))
				}
			}
		}
		// if serviceCurrentlyRunning(trip.Service, currentTime) {
		if tripCurrentlyRunning(trip, currentTime) {
			for _, stopTime := range trip.StopTimes {
				for _, foundStop := range foundStops {
					if strings.Contains(stopTime.Stop.Id, foundStop.Id) {
						stopTime.Trip = &trip
						extendedRoute := service.ERoutes[GtfsRouteID(stopTime.Trip.Route.Id)]
						// log.Debug().Msg(fmt.Sprint("rt trip: ", rtTrip))
						extendedStopTime := ExtendedStopTime{StopTime: stopTime}
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

func GetData(env *env.Env) {
	staticData, err := getStaticData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	realtimeData, err := getRealtimeData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	log.Debug().Msg("doing stuff with realtime data")
	for _, trip := range realtimeData.Trips {
		fmt.Println(len(trip.StopTimeUpdates))
	}
	currentTime := time.Now()
	fmt.Printf("VBB has %d routes and %d stations\n", len(staticData.Routes), len(staticData.Stops))
	foundStops, _ := FindStop(env, "Samariterstr.")

	var stopTimeList []gtfs.ScheduledStopTime
	for _, trip := range staticData.Trips {
		for _, stopTime := range trip.StopTimes {
			stopTime.Trip = &trip
			for _, stop := range foundStops {
				if strings.Contains(stopTime.Stop.Id, stop.Id) {
					if serviceCurrentlyRunning(stopTime.Trip.Service, currentTime) {
						stopTimeList = append(stopTimeList, stopTime)
					}
				}
			}
		}
	}
	for _, stopTime := range sortStopTimes(stopTimeList) {
		departureDuration := stopTime.DepartureTime
		departureTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(departureDuration)
		fmt.Println(departureTime, stopTime.Stop.Name, "(", stopTime.Stop.Id, ")", stopTime.Trip.Route.ShortName, "to", stopTime.Trip.Headsign, "start date:", stopTime.Trip.Service.StartDate, "end date:", stopTime.Trip.Service.EndDate)
	}
}

func getStaticData(env *env.Env) (*gtfs.Static, error) {
	if env.GtfsStaticData != nil {
		return env.GtfsStaticData, nil
	}
	gtfsSource, err := db.GetGtfsDatasource(env)
	if err != nil {
		log.Err(err).Msg("failed getting gtfs datasource")
		return &gtfs.Static{}, err
	}
	staticGtfsPath := env.App.Storage().RootURI().Path() + "staticGtfs.zip"
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
	realtimeGtfsPath := env.App.Storage().RootURI().Path() + "realtimeGtfs.bin"
	if _, err := os.Stat(realtimeGtfsPath); errors.Is(err, os.ErrNotExist) {
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

// could possibly be optimized further with early returns
func tripCurrentlyRunning(trip *gtfs.ScheduledTrip, currentTime time.Time) bool {
	// first departure time of RB69 at 13:15
	// it's currently the 03.01 at 14:30
	// subtracting 14:30 - 13:15 gives 01:15, which is still in the current day
	// therefore it can be checked if the current weekday and the service weekday are identical

	// first departure time of RB420 at 23:30
	// last departure time of RB420 at 24:30
	// it's currently the 05.02 at 00:15
	// subtracting 05.02/00:15 - 24:30 gives
	// FUCK
	// this is fucking complicated
	// maybe this entire architecture needs to be rethought

	formattedArrivalTime := currentTime.Add(-trip.StopTimes[len(trip.StopTimes)-1].ArrivalTime)
	formattedDepartureTime := currentTime.Add(-trip.StopTimes[0].DepartureTime)

	for _, addedDate := range service.AddedDates {
		if addedDate.Year() == currentTime.Year() && addedDate.Month() == currentTime.Month() && addedDate.Day() == currentTime.Day() {
			isActiveService = true
			return true
		}
	}
	for _, removedDate := range service.RemovedDates {
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

func sortStopTimes(stopTimes []gtfs.ScheduledStopTime) []gtfs.ScheduledStopTime {
	slices.SortFunc(stopTimes, func(a, b gtfs.ScheduledStopTime) int {
		return int(a.DepartureTime - b.DepartureTime)
	})
	return stopTimes
}

func SearchStopList(searchTerm string, env *env.Env) fuzzy.Ranks {
	staticData, err := getStaticData(env)
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
