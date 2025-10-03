package debugView

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/env"
	gtfsHelpers "github.com/jkulzer/transit-tool/gtfs"
	// "github.com/jamespfennell/gtfs"
)

type DebugViewWidget struct {
	widget.BaseWidget
	content fyne.CanvasObject
	env     *env.Env
}

func NewDebugViewWidget(env *env.Env) *DebugViewWidget {
	w := &DebugViewWidget{}
	w.ExtendBaseWidget(w)

	staticGtfs, err := gtfsHelpers.GetStaticData(env)
	if err != nil {
		dialog.ShowError(err, env.Window)
	}

	w.content = container.NewAppTabs(
		container.NewTabItem("Stations", NewStationSearchWidget(env, staticGtfs)),
		container.NewTabItem("Routes", NewRouteSearchWidget(env, staticGtfs)),
		container.NewTabItem("Scheduled Trips", NewStaticTripSearchWidget(env, staticGtfs)),
	)

	return w
}

func (w *DebugViewWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
