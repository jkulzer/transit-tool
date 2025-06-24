package widgets

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jkulzer/transit-tool/db"
	"github.com/jkulzer/transit-tool/env"
)

type FirstAppRunWidget struct {
	widget.BaseWidget
	content    fyne.CanvasObject
	env        *env.Env
	doneButton *widget.Button
}

func NewFirstAppRunWidget(env *env.Env) *FirstAppRunWidget {
	log.Trace().Msg("Creating Widget")
	w := &FirstAppRunWidget{}
	w.ExtendBaseWidget(w)

	w.env = env

	createGtfsSourceWidget := NewCreateGtfsSourceWidget(env)

	w.doneButton = widget.NewButton("Finish Setup", func() {
		config, err := db.GetConfig(env)
		if err != nil {
			log.Err(err).Msg("couldn't fetch config")
			dialog.ShowError(err, env.Window)
			return
		}
		config.CompletedSetup = true
		result := env.DB.Save(&config)
		if result.Error != nil {
			log.Err(err).Msg("failed to save config")
			dialog.ShowError(err, env.Window)
			return
		}
		env.Window.SetContent(NewDefaultBorderWidget(NewMainPageWidget(env)))
	})
	w.doneButton.Disable()

	w.content = container.NewVBox(
		widget.NewLabel("1. Create GTFS Datasource"),
		createGtfsSourceWidget,
		widget.NewLabel("2. Download OpenStreetMap Data"),
		widget.NewLabel("TODO"),
		w.doneButton,
	)

	return w
}

func (w *FirstAppRunWidget) CreateRenderer() fyne.WidgetRenderer {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				completedSetupSteps, err := db.SetupRequirementsFulfilled(w.env)
				// only allows the button to work when the required steps have been done
				if err != nil {
					log.Err(err).Msg("failed querying db to check if setup steps have been completed")
				} else {
					if completedSetupSteps == false {
						log.Trace().Msg("setup steps have not been completed")
						fyne.Do(func() {
							w.doneButton.Disable()
						})
					} else {
						log.Trace().Msg("setup steps have been completed")
						fyne.Do(func() {
							w.doneButton.Enable()
						})
					}
				}
			}
		}
		defer ticker.Stop()
	}()
	return widget.NewSimpleRenderer(w.content)
}
