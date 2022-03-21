package gomus_test

import (
	"testing"

	gomus "git.cesium.pw/niku/gomus/pkg"
)

func TestMapBetween(t *testing.T) {
	r := gomus.MapFloatBetween(2, 1, 3, 5, 10)
	if r != 7.5 {
		t.Fatalf("Expected value to be mapped to #, got: %f", r)
	}
}
