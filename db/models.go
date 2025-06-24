package db

import (
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	ID                uint
	CompletedSetup    bool
	DefaultGtfsSource GtfsSource `gorm:"foreignKey:ID"`
}

type GtfsSource struct {
	gorm.Model
	ID          uint
	Name        string
	StaticUrl   string
	RealtimeUrl string
}
