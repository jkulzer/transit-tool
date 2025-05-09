package env

import (
	"fyne.io/fyne/v2"
	"gorm.io/gorm"
)

type Env struct {
	DB     *gorm.DB
	App    fyne.App
	Window fyne.Window
}
