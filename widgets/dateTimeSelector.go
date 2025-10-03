package widgets

import (
	"fmt"
	// "math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jkulzer/transit-tool/env"
	"github.com/jkulzer/transit-tool/helpers"
)

type DateTimeSelectorWidget struct {
	widget.BaseWidget
	content      fyne.CanvasObject
	selectedDate time.Time
	dateText     *widget.Label
	hourWheel    *scrollLabel
	minuteWheel  *scrollLabel
}

func NewDateTimeSelectorWidget(env *env.Env) *DateTimeSelectorWidget {
	w := &DateTimeSelectorWidget{}
	w.ExtendBaseWidget(w)

	// w.hourWheel = widget.NewLabel("")
	// w.minuteWheel = widget.NewLabel("")
	w.hourWheel = newScrollLabel("", w)
	w.minuteWheel = newScrollLabel("", w)
	timeSwitcher := container.NewHBox(
		w.hourWheel,
		widget.NewLabel(":"),
		w.minuteWheel,
	)

	previousDate := widget.NewButton("<", func() {
		// subtracts one day
		w.updateTime(w.selectedDate.AddDate(0, 0, -1))
	})
	nextDate := widget.NewButton(">", func() {
		// adds one day
		w.updateTime(w.selectedDate.AddDate(0, 0, 1))
	})

	w.dateText = widget.NewLabel("")

	dateSwitcher := container.NewHBox(
		previousDate,
		w.dateText,
		nextDate,
	)

	w.content = container.NewVBox(
		timeSwitcher,
		dateSwitcher,
	)
	w.updateTime(time.Now())
	return w
}

func (w *DateTimeSelectorWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}

func (w *DateTimeSelectorWidget) updateTime(newDate time.Time) {
	w.selectedDate = newDate
	formattedDate := newDate.Format("Mon 02.01.2006")
	currentTime := time.Now()
	var dateInRelation string
	if helpers.DateIsToday(currentTime, newDate) {
		dateInRelation = "Today, "
	} else if helpers.DateIsTomorrow(currentTime, newDate) {
		dateInRelation = "Tomorrow, "
	} else if helpers.DateIsYesterday(currentTime, newDate) {
		dateInRelation = "Yesterday, "
	}
	w.dateText.SetText(dateInRelation + formattedDate)
	w.hourWheel.SetText(fmt.Sprint(newDate.Hour()))
	w.minuteWheel.SetText(fmt.Sprint(newDate.Minute()))
	w.dateText.Refresh()
}

type scrollLabel struct {
	widget.Label
	onChange               func(delta int)
	startY                 float32
	dragPosition           float32
	dateTimeSelectorWidget *DateTimeSelectorWidget
}

func newScrollLabel(initial string, dateTimeSelectorWidget *DateTimeSelectorWidget) *scrollLabel {
	sl := &scrollLabel{dateTimeSelectorWidget: dateTimeSelectorWidget}
	sl.ExtendBaseWidget(sl)
	sl.SetText(initial)
	return sl
}

// Implements fyne.Draggable
func (s *scrollLabel) Dragged(ev *fyne.DragEvent) {
	timeToAdd := time.Hour * time.Duration(ev.Dragged.DY*0.1)
	newTime := s.dateTimeSelectorWidget.selectedDate.Add(timeToAdd)
	s.dateTimeSelectorWidget.updateTime(newTime)
	log.Debug().Msg(fmt.Sprint(s.dragPosition))
}

func (s *scrollLabel) DragEnd() {}
