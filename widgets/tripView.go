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

	"github.com/jkulzer/transit-tool/colors"
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
		var delayColor color.Color
		var isDelayed bool
		var delay time.Duration
		hasRealtimeData := false

		for _, realtimeStopTime := range realtimeTrip.StopTimeUpdates {
			if *realtimeStopTime.StopID == scheduledStopTime.Stop.Id {
				hasRealtimeData = true
				isDelayed, delay = gtfsHelpers.ProcessStopTimeUpdate(realtimeStopTime, scheduledStopTime, time.Now())
				scheduleRelationship = realtimeStopTime.ScheduleRelationship
				continue
			}
		}

		var delayString string
		if hasRealtimeData {
			if isDelayed {
				delayColor = colors.Red()
			} else {
				delayColor = colors.Green()
			}
			delayString = delay.String()
		} else {
			delayColor = color.White
			delayString = ""
		}

		stopList.Add(container.NewHBox(
			canvas.NewText(gtfsHelpers.DurationToTime(scheduledStopTime.DepartureTime).Format("15:04:05"), color.White),
			canvas.NewText("delay: "+delayString, delayColor),
			canvas.NewText(scheduledStopTime.Stop.Name, color.White),
			widget.NewLabel(fmt.Sprint(scheduleRelationship)),
		))
	}
	var weekdaysActive []string
	if scheduledTrip.Service.Monday {
		weekdaysActive = append(weekdaysActive, "Monday")
	}
	if scheduledTrip.Service.Tuesday {
		weekdaysActive = append(weekdaysActive, "Tuesday")
	}
	if scheduledTrip.Service.Wednesday {
		weekdaysActive = append(weekdaysActive, "Wednesday")
	}
	if scheduledTrip.Service.Thursday {
		weekdaysActive = append(weekdaysActive, "Thursday")
	}
	if scheduledTrip.Service.Friday {
		weekdaysActive = append(weekdaysActive, "Friday")
	}
	if scheduledTrip.Service.Saturday {
		weekdaysActive = append(weekdaysActive, "Saturday")
	}
	if scheduledTrip.Service.Sunday {
		weekdaysActive = append(weekdaysActive, "Sunday")
	}

	w.content = container.NewVScroll(
		container.NewBorder(
			container.NewVBox(
				widget.NewLabel("route ID: "+scheduledTrip.Route.Id),
				widget.NewLabel("scheduled trip ID: "+scheduledTrip.ID),
				widget.NewLabel("realtime trip ID: "+realtimeTrip.ID.ID),
				widget.NewLabel("service is active on:  "+fmt.Sprint(weekdaysActive)),
			),
			nil, nil, nil,
			stopList,
		),
	)
	return w
}

func (w *TripViewWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
