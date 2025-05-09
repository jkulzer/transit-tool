package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"github.com/jkulzer/transit-tool/db"
	"github.com/jkulzer/transit-tool/env"
)

type CreateGtfsSourceWidget struct {
	widget.BaseWidget
	content      fyne.CanvasObject
	env          *env.Env
	createButton *widget.Button
	name         *widget.Entry
	url          *widget.Entry
}

func NewCreateGtfsSourceWidget(env *env.Env) *CreateGtfsSourceWidget {
	log.Trace().Msg("Creating GTFS Source Widget")
	w := &CreateGtfsSourceWidget{}
	w.ExtendBaseWidget(w)

	w.env = env

	w.name = widget.NewEntry()
	// the regex allows everything that starts and ends with a letter
	w.name.Validator = validation.NewRegexp("^[A-z].*[A-z]$", "Name must start and end with a letter")
	w.name.SetPlaceHolder("VBB GTFS Source")
	w.url = widget.NewEntry()
	w.url.Validator = validation.NewRegexp(
		// the regex should verify a http/https url
		"^(http|https)://[a-zA-Z0-9\\-\\.]+\\.[a-zA-Z]{2,}(/\\S*)?$",
		"Source must be valid URL",
	)
	w.url.SetPlaceHolder("https://production.gtfsrt.vbb.de/data")

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Name", Widget: w.name},
			{Text: "URL", Widget: w.url},
		},
		SubmitText: "Create",
		OnSubmit: func() {
			err := w.name.Validate()
			if err != nil {
				dialog.ShowError(err, env.Window)
				return
			}
			err = w.url.Validate()
			if err != nil {
				dialog.ShowError(err, env.Window)
				return
			}
			datasource := db.GtfsSource{}
			result := w.env.DB.Create(&datasource)
			if result.Error != nil {
				log.Panic().Msg("Failed to create GTFS datasource")
				dialog.ShowError(result.Error, env.Window)
			} else {
				log.Info().Msg("Submitted new GTFS datasource")
			}
		},
	}
	w.content = container.NewVBox(form)

	return w
}

func (w *CreateGtfsSourceWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
