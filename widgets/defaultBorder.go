package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jkulzer/transit-tool/env"
)

type DefaultBorderWidget struct {
	widget.BaseWidget
	content *fyne.Container
}

func NewDefaultBorderWidget(env *env.Env, center fyne.CanvasObject) *DefaultBorderWidget {
	w := &DefaultBorderWidget{}
	w.content = container.NewBorder(
		widget.NewToolbar(
			widget.NewToolbarAction(theme.MenuIcon(), func() {
				log.Debug().Msg("This should hopefully open a menu")
				settingsDialog := dialog.NewCustom("Settings", "Close", NewSettingsWidget(env), env.Window)
				settingsDialog.Show()
			}),
		),
		nil, nil, nil, center,
	)

	return w
}

func (w *DefaultBorderWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
