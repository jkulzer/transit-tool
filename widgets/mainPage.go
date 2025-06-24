package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/env"
	"github.com/jkulzer/transit-tool/gtfs"
)

type MainPageWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewMainPageWidget(env *env.Env) *MainPageWidget {
	w := &MainPageWidget{}
	w.content = container.NewVBox(
		widget.NewLabel("Main Page"),
	)

	gtfs.GetData(env)

	return w
}

func (w *MainPageWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
