package db

import (
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	CompletedSetup    bool
	DefaultGtfsSource GtfsSource `gorm:"foreignKey:ID"`
}

type GtfsSource struct {
	ID         uint
	DataPath   string
	SourceHash string
}
