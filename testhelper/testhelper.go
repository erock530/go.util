package testhelper

import (
	"math"
	"time"
)

//AlmostEqf32 tests if 2 32bit floating point numbers are within an epsilon value
func AlmostEqf32(f1 float32, f2 float32, ep float32) bool {
	return AlmostEqf64(float64(f1), float64(f2), float64(ep))
}

//AlmostEqf64 tests if 2 64bit floating point numbers are within an epsilon value
func AlmostEqf64(f1 float64, f2 float64, ep float64) bool {
	return (math.Abs(f1-f2) < ep) && (math.Abs(f2-f1) < ep)
}

//AlmostEqTime tests if time.Time numbers are within an epsilon value
func AlmostEqTime(f1 time.Time, f2 time.Time, ep time.Duration) bool {
	if (f1.Sub(f2) < ep) && (f2.Sub(f1) < ep) {
		return true
	}
	return false
}
