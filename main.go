package main

import (
	// "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jkulzer/transit-tool/db"
	"github.com/jkulzer/transit-tool/widgets"

	"github.com/rs/zerolog/log"

	"fmt"
)

func main() {
	fmt.Println("Data from:")
	fmt.Println("© OpenStreetMap contributors: https://openstreetmap.org/copyright")

	a := app.NewWithID("dev.jkulzer.transit-tool")
	w := a.NewWindow("Platform Routing App")

	dbConn := db.Init(a)

	completedSetup, err := db.HasCompletedSetup(dbConn)
	if err != nil {
		log.Warn().Msg("failed to get info if is setup process has been completed")
	}

	if completedSetup {
		center := widget.NewLabel("TODO")

		w.SetContent(center)
	} else {
		center := container.NewVBox(widgets.NewCreateGtfsSourceWidget(dbConn))
		w.SetContent(center)
	}
	w.ShowAndRun()
}
