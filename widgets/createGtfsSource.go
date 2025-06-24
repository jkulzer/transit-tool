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
	staticUrl    *widget.Entry
	realtimeUrl  *widget.Entry
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
	w.staticUrl = widget.NewEntry()
	w.staticUrl.Validator = validation.NewRegexp(
		// the regex should verify a http/https url
		"^(http|https)://[a-zA-Z0-9\\-\\.]+\\.[a-zA-Z]{2,}(/\\S*)?$",
		"Source must be valid URL",
	)
	w.staticUrl.SetPlaceHolder("https://www.vbb.de/fileadmin/user_upload/VBB/Dokumente/API-Datensaetze/gtfs-mastscharf/GTFS.zip")
	w.staticUrl.Text = "https://www.vbb.de/fileadmin/user_upload/VBB/Dokumente/API-Datensaetze/gtfs-mastscharf/GTFS.zip"

	w.realtimeUrl = widget.NewEntry()
	w.realtimeUrl.Validator = validation.NewRegexp(
		// the regex should verify a http/https url
		"^(http|https)://[a-zA-Z0-9\\-\\.]+\\.[a-zA-Z]{2,}(/\\S*)?$",
		"Source must be valid URL",
	)
	w.realtimeUrl.SetPlaceHolder("https://production.gtfsrt.vbb.de/data")
	w.realtimeUrl.Text = "https://production.gtfsrt.vbb.de/data"

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Name", Widget: w.name},
			{Text: "Static GTFS URL", Widget: w.staticUrl},
			{Text: "Realtime GTFS URL", Widget: w.realtimeUrl},
		},
		SubmitText: "Create",
		OnSubmit: func() {
			err := w.name.Validate()
			if err != nil {
				dialog.ShowError(err, env.Window)
				return
			}
			err = w.staticUrl.Validate()
			if err != nil {
				dialog.ShowError(err, env.Window)
				return
			}
			err = w.realtimeUrl.Validate()
			if err != nil {
				dialog.ShowError(err, env.Window)
				return
			}
			datasource, err := db.GetGtfsDatasource(env)
			datasource.Name = w.name.Text
			datasource.StaticUrl = w.staticUrl.Text
			datasource.RealtimeUrl = w.realtimeUrl.Text
			result := w.env.DB.Save(&datasource)
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
