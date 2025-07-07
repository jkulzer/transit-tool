package widgets

import (
	"image/color"

	"github.com/rs/zerolog/log"

	"github.com/jamespfennell/gtfs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	// "github.com/jkulzer/transit-tool/colors"
	"github.com/jkulzer/transit-tool/env"
)

type TripViewWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewTripViewWidget(env *env.Env, scheduledTrip gtfs.ScheduledTrip, realtimeTrip gtfs.Trip) *TripViewWidget {
	log.Debug().Msg("creating trip view widget")
	w := &TripViewWidget{}
	w.ExtendBaseWidget(w)

	w.content = container.NewVBox()

	for _, scheduledStopTime := range scheduledTrip.StopTimes {

		// for _, realtimeStopTime := range realtimeTrip.StopTimeUpdates {
		// 	if *realtimeStopTime.StopID == scheduledStopTime.Stop.Id {
		//
		// 	}
		// }

		w.content.Add(container.NewHBox(
			canvas.NewText(scheduledStopTime.DepartureTime.String(), color.White),
			canvas.NewText(scheduledStopTime.Stop.Name, color.White),
		))
	}
	return w
}

func (w *TripViewWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
