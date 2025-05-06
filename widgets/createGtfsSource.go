package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

type CreateGtfsSourceWidget struct {
	widget.BaseWidget
	content fyne.CanvasObject
	db      *gorm.DB
	test    string
}

func NewCreateGtfsSourceWidget(dbConn *gorm.DB) *CreateGtfsSourceWidget {
	w := &CreateGtfsSourceWidget{}
	w.ExtendBaseWidget(w)

	w.db = dbConn

	name := widget.NewEntry()
	name.Text = "Name"
	url := widget.NewEntry()
	url.Text = "URL"

	button := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		log.Info().Msg("tapped home")
	})

	w.content = container.NewVBox(name, url, button)

	return w
}

func (w *CreateGtfsSourceWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
