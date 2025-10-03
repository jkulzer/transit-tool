package debugView

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jamespfennell/gtfs"

	"github.com/jkulzer/transit-tool/env"
)

type RouteSearchWidget struct {
	widget.BaseWidget
	content fyne.CanvasObject
	env     *env.Env
}

func NewRouteSearchWidget(env *env.Env, staticGtfs *gtfs.Static) *RouteSearchWidget {
	w := &RouteSearchWidget{env: env}
	w.ExtendBaseWidget(w)

	results := container.NewVBox()
	resultsScroll := container.NewVScroll(results)

	input := widget.NewEntry()
	input.SetPlaceHolder("Route name")
	input.OnChanged = func(currentInput string) {
		results.RemoveAll()
		if len(currentInput) > 0 {
			for _, route := range staticGtfs.Routes {
				if strings.Contains(route.ShortName, currentInput) || strings.Contains(route.LongName, currentInput) {
					name := "Route has no name"
					if route.LongName != "" {
						name = route.LongName
					} else {
						if route.ShortName != "" {
							name = route.ShortName
						}
					}
					results.Add(
						container.NewHBox(widget.NewLabel(name), widget.NewLabel(route.Id)),
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

func (w *RouteSearchWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
