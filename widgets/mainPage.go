package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/env"
)

type MainPageWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewMainPageWidget(env *env.Env) *MainPageWidget {
	w := &MainPageWidget{}
	w.ExtendBaseWidget(w)

	tabs := container.NewAppTabs(
		container.NewTabItem("Departure", NewDepartureSearchWidget(env)),
		container.NewTabItem("Route", NewRouteSearchWidget(env)),
	)
	w.content = container.NewVBox(tabs)

	return w
}

func (w *MainPageWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
