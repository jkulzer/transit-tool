package helpers

import (
	"fmt"
	"errors"
	"strconv"
)

func ColorFromString(color string) (int64, int64, int64, error) {
	// TODO: add support for CSS colors and 3 digit hex codes
	// see https://wiki.openstreetmap.org/wiki/Key:colour for details
	if len(color) == 7 && string(color[0]) == ("#") {
		red, err := strconv.ParseInt(color[1:3], 16, 16)
		if err != nil {
			return 0, 0, 0, errors.New("failed decoding color with string " + color)
		}
		green, err := strconv.ParseInt(color[3:5], 16, 16)
		if err != nil {
			return 0, 0, 0, errors.New("failed decoding color with string " + color)
		}
		blue, err := strconv.ParseInt(color[5:7], 16, 16)
		if err != nil {
			return 0, 0, 0, errors.New("failed decoding color with string " + color)
		}
		return red, green, blue, nil
	} else if len(color) == 6 {
		red, err := strconv.ParseInt(color[0:2], 16, 16)
		if err != nil {
			return 0, 0, 0, errors.New("failed decoding color with string " + color)
		}
		green, err := strconv.ParseInt(color[2:4], 16, 16)
		if err != nil {
			return 0, 0, 0, errors.New("failed decoding color with string " + color)
		}
		blue, err := strconv.ParseInt(color[4:6], 16, 16)
		if err != nil {
			return 0, 0, 0, errors.New("failed decoding color with string " + color)
		}
		return red, green, blue, nil
	} else {
		// fmt.Println(printData)
		return 0, 0, 0, errors.New("failed decoding color, length of string is " + fmt.Sprint(len(color)))
	}
}
