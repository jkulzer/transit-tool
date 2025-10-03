package debugView

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jamespfennell/gtfs"

	"github.com/jkulzer/transit-tool/env"
)

type StationSearchWidget struct {
	widget.BaseWidget
	content fyne.CanvasObject
	env     *env.Env
}

func NewStationSearchWidget(env *env.Env, staticGtfs *gtfs.Static) *StationSearchWidget {
	w := &StationSearchWidget{env: env}
	w.ExtendBaseWidget(w)

	results := container.NewVBox()
	resultsScroll := container.NewVScroll(results)

	input := widget.NewEntry()
	input.SetPlaceHolder("Station name")
	input.OnChanged = func(currentInput string) {
		results.RemoveAll()
		if len(currentInput) >= 3 {
			for _, stop := range staticGtfs.Stops {
				if strings.Contains(stop.Name, currentInput) && stop.Parent == nil {
					results.Add(
						container.NewHBox(widget.NewLabel(stop.Name)),
					)
				}
			}
		}
	}

	w.content = container.NewBorder(
		container.NewVBox(
			input,
		),
		nil, nil, nil,
		resultsScroll,
	)
	return w
}

func (w *StationSearchWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
