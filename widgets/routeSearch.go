package widgets

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jkulzer/transit-tool/completion"
	"github.com/jkulzer/transit-tool/env"
	gtfsHelpers "github.com/jkulzer/transit-tool/gtfs"

	"github.com/jamespfennell/gtfs"
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
		route := gtfsHelpers.CalculateJourney(env, time.Now(), departureInput.Text, arrivalInput.Text)

		resultBox.Objects = nil

		resultBox.Add(widget.NewLabel("time to reach destination is " + route.Length.String()))
		for _, routeStop := range route.MemberStops {
			resultBox.Add(widget.NewLabel("through stop name: " + routeStop.Name))
			log.Debug().Msg("added through stop name label")
		}
		for _, routeTrip := range route.MemberTrips {
			log.Debug().Msg("added trip view widget for trip with id " + routeTrip.ID)
			resultBox.Add(NewTripViewWidget(env, routeTrip, gtfs.Trip{}))
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
