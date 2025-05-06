package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"gorm.io/gorm"
)

type InitialSetupWidget struct {
	widget.BaseWidget
	content fyne.CanvasObject
}

func NewInitialSetupWidget(db *gorm.DB) *InitialSetupWidget {
	item := &InitialSetupWidget{
		content: NewCreateGtfsSourceWidget(db),
	}

	return item
}

func (w *InitialSetupWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
