package widgets

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jamespfennell/gtfs"
	"github.com/jkulzer/transit-tool/completion"
	"github.com/jkulzer/transit-tool/env"
	gtfsHelpers "github.com/jkulzer/transit-tool/gtfs"
)

type RouteSearchWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewRouteSearchWidget(env *env.Env) *RouteSearchWidget {
	w := &RouteSearchWidget{}
	w.ExtendBaseWidget(w)

	resultBox := container.NewVBox()

	departureInput := completion.NewCompletionEntry([]string{})

	departureInput.OnChanged = func(searchTerm string) {
		stopResultsRanked := gtfsHelpers.SearchStopList(searchTerm, env)

		// no results
		if len(stopResultsRanked) == 0 {
			departureInput.HideCompletion()
			return
		}

		var resultList []string
		for _, result := range stopResultsRanked {
			resultList = append(resultList, result.Target)
		}

		// then show them
		departureInput.SetOptions(resultList)
		departureInput.ShowCompletion()
	}
	departureInput.SetPlaceHolder("departure")

	arrivalInput := completion.NewCompletionEntry([]string{})

	arrivalInput.OnChanged = func(searchTerm string) {
		stopResultsRanked := gtfsHelpers.SearchStopList(searchTerm, env)

		// no results
		if len(stopResultsRanked) == 0 {
			arrivalInput.HideCompletion()
			return
		}

		var resultList []string
		for _, result := range stopResultsRanked {
			resultList = append(resultList, result.Target)
		}

		// then show them
		arrivalInput.SetOptions(resultList)
		arrivalInput.ShowCompletion()
	}
	arrivalInput.SetPlaceHolder("destination")

	searchButton := widget.NewButton("Search", func() {

		// remove time.Now()
		route := gtfsHelpers.CalculateJourney(env, time.Now(), departureInput.Text, arrivalInput.Text, 3)

		resultBox.Objects = nil

		if route.Length != 0 {
			resultBox.Add(widget.NewLabel("time to reach destination is " + route.Length.String()))
		}
		for _, routeStop := range route.MemberStops {
			resultBox.Add(widget.NewLabel("through stop name: " + routeStop.Name))
			log.Debug().Msg("added through stop name label")
		}

		var departureStopTime gtfs.ScheduledStopTime
		for _, firstTripStopTime := range route.MemberTrips[0].StopTimes {
			if firstTripStopTime.Stop.Name == departureInput.Text {
				departureStopTime = firstTripStopTime
				break
			}
		}
		for memberIndex, stopTime := range route.MemberStopTimes {

			foundTrip := route.MemberTrips[memberIndex]

			paddingFactor := float32(10)
			tripView := NewTripViewWidget(env, foundTrip, gtfs.Trip{})
			tripDialog := dialog.NewCustom("Trip", "Close", tripView, env.Window)
			windowSize := env.Window.Canvas().Size()
			tripDialog.Resize(fyne.NewSize(windowSize.Height-windowSize.Height/paddingFactor, windowSize.Width-windowSize.Width/paddingFactor))
			tripDialog.Show()

			resultBox.Add(
				container.NewVBox(
					widget.NewLabel("trip headsign "+fmt.Sprint(route.MemberTrips[memberIndex].Headsign)),
					widget.NewLabel("trip id "+fmt.Sprint(foundTrip.ID)),
					widget.NewLabel("trip route name "+fmt.Sprint(foundTrip.Route.ShortName)),
					// NewTripViewWidget(env, foundTrip, gtfs.Trip{}),
					widget.NewLabel("stop "+stopTime.Stop.Name),
					widget.NewLabel("departure time "+gtfsHelpers.DurationToTime(departureStopTime.DepartureTime).String()),
					widget.NewLabel("arrival time "+gtfsHelpers.DurationToTime(stopTime.ArrivalTime).String()),
				),
			)
		}

		log.Debug().Msg("finished departure search")
	})
	scrollContainer := container.NewVScroll(resultBox)
	scrollContainer.SetMinSize(fyne.NewSize(0, 500))

	w.content = container.NewVBox(
		widget.NewLabel("Search for departures:"),
		departureInput,
		arrivalInput,
		searchButton,
		container.NewStack(scrollContainer),
	)

	return w
}

func (w *RouteSearchWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
