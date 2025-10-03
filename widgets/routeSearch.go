package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jkulzer/transit-tool/completion"
	"github.com/jkulzer/transit-tool/env"
	"github.com/jkulzer/transit-tool/gtfs"
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
		stopResultsRanked := gtfs.SearchStopList(searchTerm, env)

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
		stopResultsRanked := gtfs.SearchStopList(searchTerm, env)

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

		resultBox.Objects = nil

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
