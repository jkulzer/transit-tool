package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/env"
	"github.com/jkulzer/transit-tool/gtfs"
)

type DepartureSearchWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewDepartureSearchWidget(env *env.Env) *DepartureSearchWidget {
	w := &DepartureSearchWidget{}
	w.ExtendBaseWidget(w)

	resultBox := container.NewVBox()

	input := widget.NewEntry()

	searchButton := widget.NewButton("Search", func() {
		go func() {
			stationService := gtfs.QueryForDeparture(env, input.Text)

			resultBox.Objects = nil

			for _, eRoute := range stationService.ERoutes {
				w.content.Add(NewDepartureChipWidget(env, eRoute))
			}
			fyne.Do(func() { w.Refresh() })
			// for _, stopTime := range stopTimes {
			// 	departureDuration := stopTime.DepartureTime
			// 	departureTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(0), int(0), int(0), int(0), currentTime.Location()).Add(departureDuration)
			//
			// 	if departureTime.After(currentTime) && entryLimit > 0 {
			// 		entryLimit--
			//
			// 		routeColor := color.RGBA{255, 255, 255, 1}
			// 		colorString := stopTime.Trip.Route.Color
			// 		red, green, blue, err := helpers.ColorFromString(colorString)
			// 		if err != nil {
			// 			log.Err(err).Msg("failed to parse color " + colorString)
			// 		}
			// 		routeColor.R = uint8(red)
			// 		routeColor.G = uint8(green)
			// 		routeColor.B = uint8(blue)
			// 		routeColor.A = uint8(255)
			//
			// 		resultBox.Add(
			// 			container.NewHBox(
			// 				canvas.NewText(departureTime.Format("15:04"), color.White),
			// 				canvas.NewText(stopTime.Trip.Route.ShortName, routeColor),
			// 				canvas.NewText(stopTime.Trip.Headsign, color.White),
			// 				canvas.NewText(fmt.Sprint("direction id:", stopTime.Trip.DirectionId), color.White),
			// 				canvas.NewText(stopTime.Trip.Route.Id, color.White),
			// 			),
			// 		)
			// 	}
			// }
		}()
	})
	scrollContainer := container.NewVScroll(resultBox)
	scrollContainer.SetMinSize(fyne.NewSize(0, 500))

	w.content = container.NewVBox(
		widget.NewLabel("Search for departures:"),
		input,
		searchButton,
		container.NewStack(scrollContainer),
	)

	return w
}

func (w *DepartureSearchWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
