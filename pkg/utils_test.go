package gomus_test

import (
	"strconv"
	"testing"

	gomus "git.cesium.pw/niku/gomus/pkg"
)

func SliceEquals[T string | int](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestMapBetween(t *testing.T) {
	r := gomus.MapFloatBetween(2, 1, 3, 5, 10)
	if r != 7.5 {
		t.Fatalf("Expected value to be mapped to #, got: %f", r)
	}
}

func TestMapList(t *testing.T) {
	e := []string{"1", "2", "3"}
	r := gomus.MapList([]int{1, 2, 3}, func(i int) string {
		return strconv.Itoa(i)
	})
	if !SliceEquals(e, r) {
		t.Fatalf("Expected list to be mapped to string %s got: %s", e, r)
	}
}
