package db

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/jkulzer/transit-tool/env"

	"github.com/rs/zerolog/log"
)

func HasCompletedSetup(env *env.Env) (bool, error) {
	config := Config{}

	result := env.DB.First(&config)
	// checks if db record exists
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Warn().Msg("couldn't find config, initializing")
			return false, nil
		} else {
			log.Panic().Msg("failed reading config from db")
			return false, result.Error
		}
	} else {
		// checks if the prerequisites are actually fulfilled
		// gtfs datasource exists
		gtfsDatasource := GtfsSource{}
		result := env.DB.First(&gtfsDatasource)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return false, result.Error
			} else {
				log.Panic().Msg("failed reading Gtfs datasources from db")
				return false, result.Error
			}
		}

		log.Trace().Msg("setup completed: " + fmt.Sprint(config.CompletedSetup))
		log.Trace().Msg("config object: " + fmt.Sprint(config))
		return config.CompletedSetup, result.Error
	}
}

func SetupRequirementsFulfilled(env *env.Env) (bool, error) {
	gtfsDatasource := GtfsSource{}
	result := env.DB.First(&gtfsDatasource)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, result.Error
		} else {
			log.Panic().Msg("failed reading Gtfs datasources from db")
			return false, result.Error
		}
	}
	return true, nil
}

func GetConfig(env *env.Env) (Config, error) {
	config := Config{ID: 1}
	result := env.DB.First(&config)
	if result.Error != nil {
		return config, result.Error
	}
	return config, nil
}

func GetGtfsDatasource(env *env.Env) (GtfsSource, error) {
	source := GtfsSource{ID: 1}
	result := env.DB.First(&source)
	if result.Error != nil {
		return source, result.Error
	}
	return source, nil
}
