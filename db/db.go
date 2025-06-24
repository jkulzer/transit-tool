package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"

	"github.com/jkulzer/transit-tool/env"

	"github.com/rs/zerolog/log"

	"fmt"
)

func Init(app fyne.App) *gorm.DB {
	dbPathUri := storage.NewFileURI(filepath.Join(app.Storage().RootURI().Path(), "main.db"))

	db, err := gorm.Open(sqlite.Open(dbPathUri.Path()), &gorm.Config{})
	if err != nil {
		log.Err(err).Msg("failed to create/open db")
	} else {
		log.Trace().Msg("managed to create/open db at " + dbPathUri.Path())
	}

	// Migrate the schema
	db.AutoMigrate(&Config{})
	db.AutoMigrate(&GtfsSource{})

	_, err = GetConfig(&env.Env{DB: db})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result := db.Save(&Config{ID: 1})
			if result.Error != nil {
				log.Err(result.Error).Msg("failed to create config")
				panic(1)
			}
		} else {
			log.Panic().Msg("failed to get config: " + fmt.Sprint(err))
		}
	}

	return db
}
