package widgets

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jkulzer/transit-tool/completion"
	"github.com/jkulzer/transit-tool/env"
	"github.com/jkulzer/transit-tool/gtfs"
)

type DepartureSearchWidget struct {
	widget.BaseWidget
	content      *fyne.Container
	selectedTime time.Time
}

func NewDepartureSearchWidget(env *env.Env) *DepartureSearchWidget {
	w := &DepartureSearchWidget{}
	w.ExtendBaseWidget(w)

	resultBox := container.NewVBox()

	input := completion.NewCompletionEntry([]string{})

	input.OnChanged = func(searchTerm string) {
		stopResultsRanked := gtfs.SearchStopList(searchTerm, env)

		// no results
		if len(stopResultsRanked) == 0 {
			input.HideCompletion()
			return
		}

		var resultList []string
		for _, result := range stopResultsRanked {
			resultList = append(resultList, result.Target)
		}

		// then show them
		input.SetOptions(resultList)
		input.ShowCompletion()
	}

	dateTimeSelector := NewDateTimeSelectorWidget(env)

	// dateTimeSelector.OnChanged = func(selectedTime *time.Time) {
	// 	log.Debug().Msg("selected time is " + fmt.Sprint(selectedTime))
	// 	w.selectedTime = *selectedTime
	// }

	showDateTimeSelector := widget.NewButton("Select time", func() {
		dateTimePopup := dialog.NewCustom("Select date/time", "Dismiss", dateTimeSelector, env.Window)
		dateTimePopup.Show()
	})

	searchButton := widget.NewButton("Search", func() {
		stationService := gtfs.QueryForDeparture(env, input.Text)

		resultBox.Objects = nil

		for _, eRoute := range stationService.ERoutes {
			resultBox.Add(NewDepartureChipWidget(env, eRoute))
		}
		log.Debug().Msg("finished departure search")
	})
	scrollContainer := container.NewVScroll(resultBox)
	scrollContainer.SetMinSize(fyne.NewSize(0, 500))

	w.content = container.NewVBox(
		widget.NewLabel("Search for departures:"),
		input,
		showDateTimeSelector,
		searchButton,
		container.NewStack(scrollContainer),
	)

	return w
}

func (w *DepartureSearchWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
