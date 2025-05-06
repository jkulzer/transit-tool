package db

import (
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

func HasCompletedSetup(db *gorm.DB) (bool, error) {
	config := Config{}

	result := db.First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Warn().Msg("couldn't find config, initializing")
			return false, nil
		} else {
			log.Err(result.Error).Msg("failed reading config from db")
			return false, result.Error
		}
	} else {
		return config.CompletedSetup, result.Error
	}
}
