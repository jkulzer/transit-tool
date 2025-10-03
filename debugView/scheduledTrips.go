package debugView

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jamespfennell/gtfs"

	"github.com/jkulzer/transit-tool/env"
)

type StaticTripSearchWidget struct {
	widget.BaseWidget
	content fyne.CanvasObject
	env     *env.Env
}

func NewStaticTripSearchWidget(env *env.Env, staticGtfs *gtfs.Static) *StaticTripSearchWidget {
	w := &StaticTripSearchWidget{env: env}
	w.ExtendBaseWidget(w)

	results := widget.NewTable(
		func() (int, int) {
			return len(staticGtfs.Trips), 1
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(staticGtfs.Trips[i.Row].Headsign + "; " + staticGtfs.Trips[i.Row].Route.ShortName)
		})
	// resultsScroll := container.NewVScroll(results)

	// input := widget.NewEntry()
	// input.SetPlaceHolder("Trip headsign name")
	// input.OnChanged = func(currentInput string) {
	// 	results.RemoveAll()
	// 	if len(currentInput) >= 2 {
	// 		for _, trip := range staticGtfs.Trips {
	// 			if trip.Route.ShortName == currentInput {
	// 				results.Add(
	// 					container.NewHBox(
	// 						widget.NewLabel(trip.Headsign),
	// 						widget.NewLabel(trip.Route.ShortName),
	// 						widget.NewLabel("Trip ID: "+trip.ID),
	// 						widget.NewLabel("Route ID: "+trip.Route.Id),
	// 					),
	// 				)
	// 			}
	// 		}
	// 	}
	// }

	w.content = container.NewBorder(
		container.NewVBox(
		// input,
		),
		nil, nil, nil,
		results,
	)
	return w
}

func (w *StaticTripSearchWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
