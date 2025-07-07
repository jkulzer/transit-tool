package widgets

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	gtfsProto "github.com/jamespfennell/gtfs/proto"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/colors"
	"github.com/jkulzer/transit-tool/env"
	gtfsHelpers "github.com/jkulzer/transit-tool/gtfs"
	"github.com/jkulzer/transit-tool/helpers"
)

type DepartureChipWidget struct {
	widget.BaseWidget
	content *fyne.Container
	env     *env.Env
}

func NewDepartureChipWidget(env *env.Env, eRoute gtfsHelpers.ExtendedRoute) *DepartureChipWidget {
	log.Debug().Msg("creating departure chip widget")
	w := &DepartureChipWidget{env: env}
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
	content      *fyne.Container
	env          *env.Env
	stopTimeList []gtfsHelpers.ExtendedStopTime
}

func NewDirectionDepartureChipWidget(env *env.Env, stopTimeList []gtfsHelpers.ExtendedStopTime) *DirectionDepartureChipWidget {
	w := &DirectionDepartureChipWidget{
		env:          env,
		stopTimeList: stopTimeList,
	}
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

			departureDelayWidget := canvas.NewText("no data", color.White)

			log.Debug().Msg("trip: " + stopTime.StopTime.Trip.Headsign)

			log.Debug().Msg(fmt.Sprint(stopTime.RTTrip.ID.ScheduleRelationship))
			for _, stopTimeUpdate := range stopTime.RTTrip.StopTimeUpdates {
				// if strings.Contains(stopTime.StopTime.Stop.Id, *stopTimeUpdate.StopID) {
				if strings.Contains(stopTime.StopTime.Stop.Id, *stopTimeUpdate.StopID) && stopTime.StopTime.Stop.Id != *stopTimeUpdate.StopID {
					// if stopTime.StopTime.Stop.Id == *stopTimeUpdate.StopID {
					var delayColor color.Color
					var delayString string
					if stopTime.RTTrip.ID.ScheduleRelationship == gtfsProto.TripDescriptor_CANCELED {
						delayString = "cancelled"
						delayColor = colors.Red()
						log.Debug().Msg("trip cancelled")
					} else {
						// if the delay is stored as a difference from the scheduled time
						if stopTimeUpdate.Departure.Delay != nil {
							if *stopTimeUpdate.Departure.Delay > 0 {
								delayColor = colors.Red()
							} else {
								delayColor = colors.Green()
							}
							delayString = stopTimeUpdate.Departure.Delay.String()

							// if the delay is stored as a time of departure/arrival
						} else if stopTimeUpdate.Departure.Time != nil {
							agencyTimezone := stopTime.StopTime.Trip.Route.Agency.Timezone
							location, err := time.LoadLocation(agencyTimezone)
							if err != nil {
								log.Err(err).Msg("failed loading location for timezone " + agencyTimezone)
								break
							}
							departureDate := stopTimeUpdate.Departure.Time
							timezoneCorrectedDate := time.Date(departureDate.Year(), departureDate.Month(), departureTime.Day(), departureTime.Hour(), departureTime.Minute(), departureTime.Second(), departureTime.Nanosecond(), location)
							differenceToScheduled := departureTime.Sub(timezoneCorrectedDate)
							if differenceToScheduled > 0 {
								delayColor = colors.Red()
							} else {
								delayColor = colors.Green()
							}
							delayString = differenceToScheduled.String()
						}
					}
					departureDelayWidget = canvas.NewText(delayString, delayColor)
				}
			}
			log.Debug().Msg("length of stop time updates for trip " + stopTime.StopTime.Trip.ID + " and route " + stopTime.StopTime.Trip.Route.Id + " is " + fmt.Sprint(len(stopTime.RTTrip.StopTimeUpdates)))

			w.content.Add(
				container.NewHBox(
					container.NewVBox(
						canvas.NewText(departureTime.Format("15:04"), color.White),
						departureDelayWidget,
					),
					canvas.NewText(stopTime.StopTime.Trip.Route.ShortName, routeColor),
					canvas.NewText(stopTime.StopTime.Trip.Headsign, color.White),
					canvas.NewText("Trip ID: "+stopTime.StopTime.Trip.ID, color.White),
					// canvas.NewText(fmt.Sprint("direction id:", stopTime.Trip.DirectionId), color.White),
					// canvas.NewText(stopTime.Trip.Route.Id, color.White),
				),
			)
			log.Debug().Msg("completed subwidget for departure chip")
			return w
		}
	}

	return w
}

func (w *DirectionDepartureChipWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}

// func (w *DirectionDepartureChipWidget) OnClick() {
// 	go func() {
// 		tripView := NewTripViewWidget(w.env, *w.stopTimeList[0].StopTime.Trip, w.stopTimeList[0].RTTrip)
// 		dialog.NewCustom("Route", "Close", tripView, w.env.Window)
// 	}()
// }
