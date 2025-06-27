package widgets

import (
	"fmt"
	"image/color"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jamespfennell/gtfs"

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

func NewDirectionDepartureChipWidget(env *env.Env, stopTimes []gtfs.ScheduledStopTime) *DirectionDepartureChipWidget {
	w := &DirectionDepartureChipWidget{}
	w.ExtendBaseWidget(w)

	w.content = container.NewHBox()

	currentTime := time.Now()

	for _, stopTime := range stopTimes {
		departureDuration := stopTime.DepartureTime
		departureTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(departureDuration)

		if departureTime.After(currentTime) {

			routeColor := color.RGBA{255, 255, 255, 1}
			colorString := stopTime.Trip.Route.Color
			red, green, blue, err := helpers.ColorFromString(colorString)
			if err != nil {
				log.Err(err).Msg("failed to parse color " + colorString)
			}
			routeColor.R = uint8(red)
			routeColor.G = uint8(green)
			routeColor.B = uint8(blue)
			routeColor.A = uint8(255)

			w.content.Add(
				container.NewHBox(
					canvas.NewText(departureTime.Format("15:04"), color.White),
					canvas.NewText(stopTime.Trip.Route.ShortName, routeColor),
					canvas.NewText(stopTime.Trip.Headsign, color.White),
					canvas.NewText(fmt.Sprint("direction id:", stopTime.Trip.DirectionId), color.White),
					canvas.NewText(stopTime.Trip.Route.Id, color.White),
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
