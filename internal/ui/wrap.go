package ui

import "sync/atomic"

var wrapWidth atomic.Int64
var stableLayout atomic.Int32

func setWrapWidth(width int) {
	if width < 0 {
		return
	}
	wrapWidth.Store(int64(width))
}

func currentWrapWidth() int {
	return int(wrapWidth.Load())
}

func setStableLayout(enabled bool) {
	if enabled {
		stableLayout.Store(1)
		return
	}
	stableLayout.Store(0)
}

func currentStableLayout() bool {
	return stableLayout.Load() == 1
}
