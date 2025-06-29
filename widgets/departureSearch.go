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
				resultBox.Add(NewDepartureChipWidget(env, eRoute))
				fyne.Do(func() { w.Refresh() })
			}
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
