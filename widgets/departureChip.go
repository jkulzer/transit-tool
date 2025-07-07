package widgets

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	// "github.com/jamespfennell/gtfs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/env"
	gtfsHelpers "github.com/jkulzer/transit-tool/gtfs"
	"github.com/jkulzer/transit-tool/helpers"
)

type DepartureChipWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewDepartureChipWidget(env *env.Env, eRoute gtfsHelpers.ExtendedRoute) *DepartureChipWidget {
	w := &DepartureChipWidget{}
	w.ExtendBaseWidget(w)

	w.content = container.NewVBox(
		NewDirectionDepartureChipWidget(env, eRoute.StopTimesDirectionTrue),
		NewDirectionDepartureChipWidget(env, eRoute.StopTimesDirectionFalse),
		NewDirectionDepartureChipWidget(env, eRoute.StopTimesNoDirection),
	)

	return w
}

func (w *DepartureChipWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}

type DirectionDepartureChipWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewDirectionDepartureChipWidget(env *env.Env, stopTimeList []gtfsHelpers.ExtendedStopTime) *DirectionDepartureChipWidget {
	w := &DirectionDepartureChipWidget{}
	w.ExtendBaseWidget(w)

	w.content = container.NewHBox()

	currentTime := time.Now()

	for _, stopTime := range stopTimeList {
		departureDuration := stopTime.StopTime.DepartureTime
		departureTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(departureDuration)

		if departureTime.After(currentTime) {

			routeColor := color.RGBA{255, 255, 255, 1}
			colorString := stopTime.StopTime.Trip.Route.Color
			red, green, blue, err := helpers.ColorFromString(colorString)
			if err != nil {
				log.Err(err).Msg("failed to parse color " + colorString)
			}
			routeColor.R = uint8(red)
			routeColor.G = uint8(green)
			routeColor.B = uint8(blue)
			routeColor.A = uint8(255)

			if routeColor.R == 0 && routeColor.G == 0 && routeColor.B == 0 {
				routeColor.R = uint8(255)
				routeColor.G = uint8(255)
				routeColor.B = uint8(255)
			}

			departureDelay := "no data"

			log.Debug().Msg("trip: " + stopTime.StopTime.Trip.Headsign)

			for _, stopTimeUpdate := range stopTime.RTTrip.StopTimeUpdates {
				if strings.Contains(stopTime.StopTime.Stop.Id, *stopTimeUpdate.StopID) {
					// if stopTimeUpdate.Departure != nil {
					if stopTimeUpdate.Departure.Time != nil {
						// departureDelay = stopTimeUpdate.Departure.Time.Format("15:04")
						agencyTimezone := stopTime.StopTime.Trip.Route.Agency.Timezone
						location, err := time.LoadLocation(agencyTimezone)
						if err != nil {
							log.Err(err).Msg("failed loading location for timezone " + agencyTimezone)
							break
						}
						departureDate := stopTimeUpdate.Departure.Time
						timezoneCorrectedDate := time.Date(departureDate.Year(), departureDate.Month(), departureTime.Day(), departureTime.Hour(), departureTime.Minute(), departureTime.Second(), departureTime.Nanosecond(), location)
						departureDelay = timezoneCorrectedDate.String()
					} else if stopTimeUpdate.Departure.Delay != nil {
						departureDelay = stopTimeUpdate.Departure.Delay.String()
					}
				}
			}
			log.Debug().Msg("length of stop time updates for trip " + stopTime.StopTime.Trip.ID + " is " + fmt.Sprint(len(stopTime.RTTrip.StopTimeUpdates)))

			w.content.Add(
				container.NewHBox(
					container.NewVBox(
						canvas.NewText(departureTime.Format("15:04"), color.White),
						canvas.NewText(departureDelay, color.White),
					),
					canvas.NewText(stopTime.StopTime.Trip.Route.ShortName, routeColor),
					canvas.NewText(stopTime.StopTime.Trip.Headsign, color.White),
					canvas.NewText("Trip ID: "+stopTime.StopTime.Trip.ID, color.White),
					// canvas.NewText(fmt.Sprint("direction id:", stopTime.Trip.DirectionId), color.White),
					// canvas.NewText(stopTime.Trip.Route.Id, color.White),
				),
			)
			return w
		}
	}

	return w
}

func (w *DirectionDepartureChipWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
