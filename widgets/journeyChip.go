package widgets

import (
	// "fmt"
	// "image/color"
	// "strings"
	// "time"

	"github.com/rs/zerolog/log"

	// gtfsProto "github.com/jamespfennell/gtfs/proto"

	"fyne.io/fyne/v2"
	// "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	// "github.com/jkulzer/transit-tool/colors"
	"github.com/jkulzer/transit-tool/env"
	gtfsHelpers "github.com/jkulzer/transit-tool/gtfs"
	// "github.com/jkulzer/transit-tool/helpers"
)

type JourneyChipWidget struct {
	widget.BaseWidget
	content *fyne.Container
	env     *env.Env
}

func NewJourneyChipWidget(env *env.Env, eRoute gtfsHelpers.Journey) *JourneyChipWidget {
	log.Debug().Msg("creating departure chip widget")
	w := &JourneyChipWidget{env: env}
	w.ExtendBaseWidget(w)

	w.content = container.NewVBox()
	return w
}

func (w *JourneyChipWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}
