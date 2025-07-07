package colors

import (
	"image/color"
)

func Green() color.Color {
	green := color.RGBA{0, 255, 0, 255}
	return green
}

func Red() color.Color {
	green := color.RGBA{255, 0, 0, 255}
	return green
}
