package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"

	"github.com/rs/zerolog/log"
)

func Init(app fyne.App) *gorm.DB {
	dbPathUri := storage.NewFileURI(filepath.Join(app.Storage().RootURI().Path(), "main.db"))

	db, err := gorm.Open(sqlite.Open(dbPathUri.Path()), &gorm.Config{})
	if err != nil {
		log.Err(err).Msg("failed to create/open db")
	}

	// Migrate the schema
	db.AutoMigrate(&Config{})

	return db
}
