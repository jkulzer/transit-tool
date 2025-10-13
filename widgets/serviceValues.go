package widgets

import (
	"fmt"

	"github.com/jamespfennell/gtfs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/env"
)

type ServiceValuesWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewServiceValuesWidget(env *env.Env, service gtfs.Service) *ServiceValuesWidget {
	w := &ServiceValuesWidget{}
	w.ExtendBaseWidget(w)

	var weekdaysActive []string
	if service.Monday {
		weekdaysActive = append(weekdaysActive, "Monday")
	}
	if service.Tuesday {
		weekdaysActive = append(weekdaysActive, "Tuesday")
	}
	if service.Wednesday {
		weekdaysActive = append(weekdaysActive, "Wednesday")
	}
	if service.Thursday {
		weekdaysActive = append(weekdaysActive, "Thursday")
	}
	if service.Friday {
		weekdaysActive = append(weekdaysActive, "Friday")
	}
	if service.Saturday {
		weekdaysActive = append(weekdaysActive, "Saturday")
	}
	if service.Sunday {
		weekdaysActive = append(weekdaysActive, "Sunday")
	}

	var addedDatesString string
	for index, addedDate := range service.AddedDates {
		addedDatesString = addedDatesString + addedDate.Format("02.01.2006")
		if addedDatesString != "" && index < len(service.AddedDates)-1 {
			addedDatesString = addedDatesString + ", "
		}
	}

	var removedDatesString string
	for index, removedDate := range service.RemovedDates {
		removedDatesString = removedDatesString + removedDate.Format("02.01.2006")
		if removedDatesString != "" && index < len(service.RemovedDates)-1 {
			removedDatesString = removedDatesString + ", "
		}
	}

	w.content = container.NewVBox()
	if weekdaysActive != nil {
		w.content.Add(widget.NewLabel("service is active on:  " + fmt.Sprint(weekdaysActive)))
	}
	w.content.Add(widget.NewLabel("service is active between:  " + service.StartDate.Format("02.01.2006") + " and " + service.EndDate.Format("02.01.2006")))
	if service.AddedDates != nil {
		addedDatesWidget := widget.NewLabel("service is additionally running on: " + addedDatesString)
		addedDatesWidget.Wrapping = fyne.TextWrapWord
		w.content.Add(addedDatesWidget)
	}
	if service.RemovedDates != nil {
		removedDatesWidget := widget.NewLabel("service is not running on: " + removedDatesString)
		removedDatesWidget.Wrapping = fyne.TextWrapWord
		w.content.Add(removedDatesWidget)
	}
	return w
}

func (w *ServiceValuesWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
