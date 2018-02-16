package channellight

import (
	"math"
)

func linearScale(val float64, zeroAt int, maxAt int) int {
	rangeLen := math.Abs(float64(maxAt - zeroAt))
	norm := math.Min(math.Max(math.Abs(val-float64(zeroAt))/rangeLen, 0), 1)
	return int(math.Floor(norm*255 + 0.5))
}

func isBetween(val float64, stopOne, stopTwo int) bool {
	useStopOne := float64(stopOne)
	useStopTwo := float64(stopTwo)
	if stopOne > stopTwo {
		if val >= useStopOne || val <= useStopTwo {
			return true
		}
	} else {
		if val >= useStopOne && val <= useStopTwo {
			return true
		}
	}
	return false
}

func multiStop(val float64, zeroOne, zeroTwo, maxOne, maxTwo int) int {
	if isBetween(val, zeroOne, zeroTwo) {
		return 0
	}
	if isBetween(val, maxOne, maxTwo) {
		return 255
	}
	if isBetween(val, zeroTwo, maxOne) {
		return linearScale(val, zeroTwo, maxOne)
	}
	return linearScale(val, zeroOne, maxTwo)
}
