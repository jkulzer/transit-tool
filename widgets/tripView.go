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
)

type TripViewWidget struct {
	widget.BaseWidget
	content *container.Scroll
}

func NewTripViewWidget(env *env.Env, scheduledTrip gtfs.ScheduledTrip, realtimeTrip gtfs.Trip) *TripViewWidget {
	log.Debug().Msg("creating trip view widget")
	w := &TripViewWidget{}
	w.ExtendBaseWidget(w)

	stopList := container.NewVBox()

	for _, scheduledStopTime := range scheduledTrip.StopTimes {
		var scheduleRelationship gtfs.StopTimeUpdateScheduleRelationship
		var arrivalDelayStatus gtfsHelpers.DelayStatus
		var arrivalDelay time.Duration
		var departureDelayStatus gtfsHelpers.DelayStatus
		var departureDelay time.Duration
		hasRealtimeData := false

		for _, realtimeStopTime := range realtimeTrip.StopTimeUpdates {
			if *realtimeStopTime.StopID == scheduledStopTime.Stop.Id {
				hasRealtimeData = true
				arrivalDelay, arrivalDelayStatus, departureDelay, departureDelayStatus = gtfsHelpers.ProcessStopTimeUpdate(realtimeStopTime, scheduledStopTime, time.Now())
				scheduleRelationship = realtimeStopTime.ScheduleRelationship
				continue
			}
		}

		arrivalDelayString, arrivalDelayColor := gtfsHelpers.GetDelayColor(arrivalDelay, arrivalDelayStatus, hasRealtimeData)
		departureDelayString, departureDelayColor := gtfsHelpers.GetDelayColor(departureDelay, departureDelayStatus, hasRealtimeData)

		arrivalTimeString := gtfsHelpers.DurationToTime(scheduledStopTime.ArrivalTime).Format("15:04:05")
		departureTimeString := gtfsHelpers.DurationToTime(scheduledStopTime.DepartureTime).Format("15:04:05")

		var stopTimeDisplay *fyne.Container
		if scheduledStopTime.ArrivalTime == scheduledStopTime.DepartureTime && arrivalDelay == departureDelay {
			stopTimeDisplay = container.NewHBox(
				canvas.NewText(departureTimeString, color.White),
				canvas.NewText(departureDelayString, departureDelayColor),
			)
		} else {
			stopTimeDisplay = container.NewHBox(
				container.NewVBox(
					canvas.NewText(arrivalTimeString, color.White),
					canvas.NewText(departureTimeString, color.White),
				),
				container.NewVBox(
					canvas.NewText(arrivalDelayString, arrivalDelayColor),
					canvas.NewText(departureDelayString, departureDelayColor),
				),
				canvas.NewText("Dwell time: "+time.Duration(int64(scheduledStopTime.DepartureTime)-int64(scheduledStopTime.ArrivalTime)).String(), color.White),
			)
		}

		stopList.Add(container.NewHBox(
			stopTimeDisplay,
			canvas.NewText(scheduledStopTime.Stop.Name, color.White),
			widget.NewLabel(fmt.Sprint(scheduleRelationship)),
		))
	}

	w.content = container.NewVScroll(
		container.NewBorder(
			// top
			nil,
			// bottom
			container.NewVBox(
				widget.NewLabel("route ID: "+scheduledTrip.Route.Id),
				widget.NewLabel("scheduled trip ID: "+scheduledTrip.ID),
				widget.NewLabel("realtime trip ID: "+realtimeTrip.ID.ID),
				NewServiceValuesWidget(env, *scheduledTrip.Service),
			),
			nil, nil,
			stopList,
		),
	)
	return w
}

func (w *TripViewWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
