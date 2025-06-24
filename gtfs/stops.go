package gtfs

import (
	"github.com/jamespfennell/gtfs"
)

// the stop.Root() function returns the same stop if the stop is the topmost stop, so this is a wrapper
func stopIsTopLevel(stop gtfs.Stop) bool {
	rootStop := stop.Root()
	if *rootStop == stop {
		return true
	} else {
		return false
	}
}
