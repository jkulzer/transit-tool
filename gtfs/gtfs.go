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
	"strings"
	"time"
)

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

func GetData(env *env.Env) {
	staticData, err := getStaticData(env)
	if err != nil {
		log.Err(err).Msg("failed to get static data")
	}
	currentTime := time.Now()
	fmt.Printf("VBB has %d routes and %d stations\n", len(staticData.Routes), len(staticData.Stops))
	foundStops, _ := FindStop(env, "Brandenburger Tor")
	for _, stop := range foundStops {
		fmt.Println(stop.Name+","+stop.Description+","+fmt.Sprint(*stop.Latitude), fmt.Sprint(*stop.Longitude), " Code: "+stop.Id)
		if stop.Id == "de:11000:900100025" {
			for _, trip := range staticData.Trips {
				for _, stopTime := range trip.StopTimes {
					if strings.Contains(stopTime.Stop.Id, stop.Id) {
						departureDuration := stopTime.DepartureTime
						departureTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(departureDuration)
						fmt.Println(departureTime, stopTime.Stop.Name, "(", stopTime.Stop.Id, ")", trip.Route.ShortName, "to", trip.Headsign)
					}
				}
			}
			for _, stopTime := range env.GtfsStaticOptimized.StopTimesByStop[stop.Id] {
				fmt.Println(*stopTime)
			}
		}
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
