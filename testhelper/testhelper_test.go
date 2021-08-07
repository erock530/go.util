package testhelper

import (
	"math"
	"testing"
)

func TestAlmostEqf64(t *testing.T) {
	f1, f2, ep := 1.0, 1.001, .01
	if !AlmostEqf64(f1, f2, ep) {
		t.Fatalf("%v and %v should be equal with ep value of %v", f1, f2, ep)
	}

	f1, f2, ep = 1.0, 1.001, .0001
	if AlmostEqf64(f1, f2, ep) {
		t.Fatalf("%v and %v should not be equal with ep value of %v", f1, f2, ep)
	}

}

func TestAlmostEqf32(t *testing.T) {
	t.Logf("%t", math.Abs(1.-1.001) < .01 && math.Abs(1.001-1.) < .01)
	f1, f2, ep := float32(1.0), float32(1.001), float32(.01)
	if !AlmostEqf32(f1, f2, ep) {
		t.Errorf("%v and %v should be equal with ep value of %v", f1, f2, ep)
	}

	f1, f2, ep = float32(1.0), float32(1.01), float32(.009)
	if AlmostEqf32(f1, f2, ep) {
		t.Errorf("%v and %v should not be equal with ep value of %v", f1, f2, ep)
	}
}
