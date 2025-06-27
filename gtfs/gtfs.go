package gtfs

import (
	"github.com/jkulzer/transit-tool/db"
	"github.com/jkulzer/transit-tool/env"

	"github.com/jamespfennell/gtfs"

	"github.com/rs/zerolog/log"

	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

type StationService struct {
	// IFOPT code
	// https://en.wikipedia.org/wiki/Identification_of_Fixed_Objects_in_Public_Transport
	// could potentially also be something different in different systems, in the VBB GTFS it's IFOPT though
	StopID  string
	ERoutes map[string]ExtendedRoute // string is gtfs route id
}

type ExtendedRoute struct {
	Route                   gtfs.Route
	StopTimesDirectionTrue  []gtfs.ScheduledStopTime
	StopTimesDirectionFalse []gtfs.ScheduledStopTime
	StopTimesNoDirection    []gtfs.ScheduledStopTime
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

func QueryForDeparture(env *env.Env, stopName string) StationService {
	currentTime := time.Now()
	staticData, err := getStaticData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	log.Debug().Msg("finished parsing data (" + fmt.Sprint(time.Since(currentTime)) + ")")
	foundStops, _ := FindStop(env, stopName)
	log.Debug().Msg("finished finding stops (" + fmt.Sprint(time.Since(currentTime)) + ")")

	service := StationService{}
	service.ERoutes = make(map[string]ExtendedRoute)

	for _, trip := range staticData.Trips {
		if serviceCurrentlyRunning(trip.Service, currentTime) {
			for _, stopTime := range trip.StopTimes {
				for _, stop := range foundStops {
					if strings.Contains(stopTime.Stop.Id, stop.Id) {
						stopTime.Trip = &trip
						extendedRoute := service.ERoutes[stopTime.Trip.Route.Id]
						switch trip.DirectionId {
						case gtfs.DirectionID_Unspecified:
							extendedRoute.StopTimesNoDirection = append(extendedRoute.StopTimesNoDirection, stopTime)
						case gtfs.DirectionID_True:
							extendedRoute.StopTimesDirectionTrue = append(extendedRoute.StopTimesDirectionTrue, stopTime)
						case gtfs.DirectionID_False:
							extendedRoute.StopTimesDirectionFalse = append(extendedRoute.StopTimesDirectionFalse, stopTime)
						}
						service.ERoutes[stopTime.Trip.Route.Id] = extendedRoute
					}
				}
			}
		}
	}
	log.Debug().Msg("finished creating list(" + fmt.Sprint(time.Since(currentTime)) + ")")
	for key, eRoute := range service.ERoutes {
		eRoute.StopTimesDirectionTrue = sortStopTimes(eRoute.StopTimesDirectionTrue)
		eRoute.StopTimesDirectionFalse = sortStopTimes(eRoute.StopTimesDirectionFalse)
		eRoute.StopTimesNoDirection = sortStopTimes(eRoute.StopTimesNoDirection)
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
func serviceCurrentlyRunning(service *gtfs.Service, currentTime time.Time) bool {
	isActiveService := currentTime.After(service.StartDate) && currentTime.Before(service.EndDate)

	switch currentTime.Weekday() {
	case time.Monday:
		if service.Monday {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Tuesday:
		if service.Tuesday {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Wednesday:
		if service.Wednesday {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Thursday:
		if service.Thursday {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Friday:
		if service.Friday {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Saturday:
		if service.Saturday {
			isActiveService = true
		} else {
			isActiveService = false
		}
	case time.Sunday:
		if service.Sunday {
			isActiveService = true
		} else {
			isActiveService = false
		}
	}

	for _, addedDate := range service.AddedDates {
		if addedDate.Year() == currentTime.Year() && addedDate.Month() == currentTime.Month() && addedDate.Day() == currentTime.Day() {
			isActiveService = true
		}
	}
	for _, removedDate := range service.RemovedDates {
		if removedDate.Year() == currentTime.Year() && removedDate.Month() == currentTime.Month() && removedDate.Day() == currentTime.Day() {
			isActiveService = false
		}
	}
	return isActiveService
}

func sortStopTimes(stopTimes []gtfs.ScheduledStopTime) []gtfs.ScheduledStopTime {
	slices.SortFunc(stopTimes, func(a, b gtfs.ScheduledStopTime) int {
		return int(a.DepartureTime - b.DepartureTime)
	})
	return stopTimes
}
