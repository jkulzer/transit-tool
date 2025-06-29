package env

import (
	"fyne.io/fyne/v2"
	"gorm.io/gorm"

	"github.com/jamespfennell/gtfs"
)

type Env struct {
	DB                  *gorm.DB
	App                 fyne.App
	Window              fyne.Window
	GtfsStaticData      *gtfs.Static
	GtfsRealtimeData    *gtfs.Realtime
	GtfsStaticOptimized GtfsStaticOptimized
}

type GtfsStaticOptimized struct {
	Stops           []*gtfs.Stop
	StopTimesByStop map[string][]*gtfs.ScheduledStopTime
}
