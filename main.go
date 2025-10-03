package main

import (
	// "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/jkulzer/transit-tool/db"
	"github.com/jkulzer/transit-tool/debugView"
	"github.com/jkulzer/transit-tool/env"
	"github.com/jkulzer/transit-tool/widgets"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"fmt"
	"os"
)

func main() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	fmt.Println("Data from:")
	fmt.Println("© OpenStreetMap contributors: https://openstreetmap.org/copyright")

	a := app.NewWithID("dev.jkulzer.transit-tool")
	w := a.NewWindow("Platform Routing App")

	dbConn := db.Init(a)

	env := env.Env{
		DB:     dbConn,
		App:    a,
		Window: w,
	}

	args := os.Args[1:]
	completedSetup, err := db.HasCompletedSetup(&env)
	if err != nil {
		log.Warn().Msg("failed to get info if is setup process has been completed")
	}

	if completedSetup {
		log.Trace().Msg("Setup is completed")
		if args[0] == "debug" {
			w.SetContent(debugView.NewDebugViewWidget(&env))
		} else {
			center := widgets.NewDefaultBorderWidget(&env, widgets.NewMainPageWidget(&env))
			w.SetContent(center)
		}
	} else {
		log.Trace().Msg("Setup is not completed")
		center := container.NewVBox(widgets.NewFirstAppRunWidget(&env))
		w.SetContent(center)
	}
	w.ShowAndRun()
}
