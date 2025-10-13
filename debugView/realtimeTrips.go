package debugView

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jamespfennell/gtfs"

	"github.com/jkulzer/transit-tool/env"
)

type RealtimeTripSearchWidget struct {
	widget.BaseWidget
	content fyne.CanvasObject
	env     *env.Env
}

func NewRealtimeTripSearchWidget(env *env.Env, realtimeGtfs *gtfs.Realtime) *RealtimeTripSearchWidget {
	w := &RealtimeTripSearchWidget{env: env}
	w.ExtendBaseWidget(w)

	results := container.NewVBox()
	resultsScroll := container.NewVScroll(results)

	input := widget.NewEntry()
	input.SetPlaceHolder("Realtime Trip ID")
	searchButton := widget.NewButton("Search", func() {
		results.RemoveAll()
		if len(input.Text) >= 2 {
			for _, trip := range realtimeGtfs.Trips {
				if trip.ID.ID == input.Text {
					results.Add(
						container.NewHBox(
							widget.NewLabel("ID: "+trip.ID.ID),
							widget.NewLabel("Route ID: "+trip.ID.RouteID),
							widget.NewLabel("Schedule Relationship: "+fmt.Sprint(trip.ID.ScheduleRelationship)),
						),
					)
				}
			}
		}
	})

	w.content = container.NewBorder(
		container.NewVBox(
			input,
			searchButton,
		),
		nil, nil, nil,
		resultsScroll,
	)
	return w
}

func (w *RealtimeTripSearchWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
