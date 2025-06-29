package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/db"
	"github.com/jkulzer/transit-tool/env"
)

type SettingsWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewSettingsWidget(env *env.Env) *SettingsWidget {
	w := &SettingsWidget{}

	config, err := db.GetConfig(env)
	if err != nil {
		dialog.ShowError(err, env.Window)
		return w
	}

	staticEntry := widget.NewEntry()
	staticEntry.Text = config.DefaultGtfsSource.StaticUrl
	realtimeEntry := widget.NewEntry()
	realtimeEntry.Text = config.DefaultGtfsSource.RealtimeUrl

	saveButton := widget.NewButton("Save", func() {

	})

	w.content = container.NewVBox(
		widget.NewLabel("Name"),
		widget.NewLabel("Static URL"),
		staticEntry,
		widget.NewLabel("Realtime URL"),
		realtimeEntry,
	)

	return w
}

func (w *SettingsWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
