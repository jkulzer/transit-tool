package widgets

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

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

		departureStationID := gtfsHelpers.StopNameToIFOPT(env, departureInput.Text)
		arrivalStationID := gtfsHelpers.StopNameToIFOPT(env, arrivalInput.Text)
		journey := gtfsHelpers.CalculateJourney(env, time.Now(), departureStationID, arrivalStationID, 3)

		resultBox.Objects = nil

		arrivalStopLabel := journey[len(journey)-1][arrivalStationID]
		for i := len(journey) - 1; i >= 0; i-- {
		}

		resultBox.Add(widget.NewLabel("arrival time: " + fmt.Sprint(journey)))

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
