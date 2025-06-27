package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/env"
)

type MainPageWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewMainPageWidget(env *env.Env) *MainPageWidget {
	w := &MainPageWidget{}
	w.ExtendBaseWidget(w)

	w.content = container.NewHBox(
		NewDepartureSearchWidget(env),
	)

	return w
}

func (w *MainPageWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
