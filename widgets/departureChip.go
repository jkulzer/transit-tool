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
	"fyne.io/fyne/v2/dialog"
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

	w.content = container.NewVBox()
	if len(eRoute.StopTimesDirectionTrue) > 0 {
		w.content.Add(NewDirectionDepartureChipWidget(env, eRoute.StopTimesDirectionTrue))
	}
	if len(eRoute.StopTimesDirectionFalse) > 0 {
		w.content.Add(NewDirectionDepartureChipWidget(env, eRoute.StopTimesDirectionFalse))
	}
	if len(eRoute.StopTimesNoDirection) > 0 {
		w.content.Add(NewDirectionDepartureChipWidget(env, eRoute.StopTimesNoDirection))
	}
	return w
}

func (w *DepartureChipWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}

type DirectionDepartureChipWidget struct {
	widget.BaseWidget
	content         *fyne.Container
	env             *env.Env
	stopTimeList    []gtfsHelpers.ExtendedStopTime
	visibleStopTime gtfsHelpers.ExtendedStopTime
}

func NewDirectionDepartureChipWidget(env *env.Env, stopTimeList []gtfsHelpers.ExtendedStopTime) *DirectionDepartureChipWidget {
	w := &DirectionDepartureChipWidget{
		env:          env,
		stopTimeList: stopTimeList,
	}
	w.ExtendBaseWidget(w)

	w.content = container.NewHBox()

	currentTime := time.Now()

	log.Debug().Msg("creating direction departure widget for route " + fmt.Sprint(stopTimeList[0].StopTime.Trip.Route.ShortName))

	for _, stopTime := range stopTimeList {
		departureDuration := stopTime.StopTime.DepartureTime
		departureTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(departureDuration)

		if departureTime.After(currentTime) && departureTime.Sub(time.Now()) < time.Hour {

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

			w.visibleStopTime = stopTime

			departureDelayWidget := canvas.NewText("no data", color.White)

			log.Debug().Msg("trip: " + stopTime.StopTime.Trip.Headsign)

			log.Debug().Msg(fmt.Sprint(stopTime.RTTrip.ID.ScheduleRelationship))
			for _, stopTimeUpdate := range stopTime.RTTrip.StopTimeUpdates {
				if stopTimeUpdate.StopID == nil || stopTime.StopTime.Stop.Id == "" {
					log.Debug().Msg(fmt.Sprint(stopTimeUpdate))
					log.Debug().Msg(fmt.Sprint(stopTime.StopTime))
				}
				if strings.Contains(stopTime.StopTime.Stop.Id, *stopTimeUpdate.StopID) && stopTime.StopTime.Stop.Id != *stopTimeUpdate.StopID {
					var delayColor color.Color
					var delayString string
					if stopTime.RTTrip.ID.ScheduleRelationship == gtfsProto.TripDescriptor_CANCELED {
						delayString = "cancelled"
						delayColor = colors.Red()
						log.Debug().Msg("trip cancelled")
					} else {
						isDelayed, delay := gtfsHelpers.ProcessStopTimeUpdate(stopTimeUpdate, stopTime.StopTime, currentTime)
						if isDelayed {
							delayColor = colors.Red()
						} else {
							delayColor = colors.Green()
						}
						delayString = delay.String()
					}
					departureDelayWidget = canvas.NewText(delayString, delayColor)
				}
			}
			log.Debug().Msg("length of stop time updates for trip " + stopTime.StopTime.Trip.ID + " and route " + stopTime.StopTime.Trip.Route.Id + " is " + fmt.Sprint(len(stopTime.RTTrip.StopTimeUpdates)))

			w.content.Add(
				container.NewHBox(
					container.NewVBox(
						canvas.NewText(departureTime.Format("15:04:05"), color.White),
						departureDelayWidget,
					),
					canvas.NewText(stopTime.StopTime.Trip.Route.ShortName, routeColor),
					canvas.NewText(stopTime.StopTime.Trip.Headsign, color.White),
					canvas.NewText("Trip ID: "+stopTime.StopTime.Trip.ID, color.White),
					canvas.NewText(fmt.Sprint("direction id:", stopTime.StopTime.Trip.DirectionId), color.White),
					canvas.NewText(fmt.Sprint("route id:", stopTime.StopTime.Trip.Route.Id), color.White),
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

func (w *DirectionDepartureChipWidget) Tapped(_ *fyne.PointEvent) {
	paddingFactor := float32(10)
	tripView := NewTripViewWidget(w.env, *w.visibleStopTime.StopTime.Trip, w.visibleStopTime.RTTrip)
	tripDialog := dialog.NewCustom("Trip", "Close", tripView, w.env.Window)
	windowSize := w.env.Window.Canvas().Size()
	tripDialog.Resize(fyne.NewSize(windowSize.Height-windowSize.Height/paddingFactor, windowSize.Width-windowSize.Width/paddingFactor))
	tripDialog.Show()
}
