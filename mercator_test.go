package main

import "testing"

func TestLatLong2XY(t *testing.T) {
	cases := []struct {
		lat  float64
		long float64
		zoom int
		x    int
		y    int
	}{
		{0, 0, 0, 0, 0},

		{50, -180, 1, 0, 0},
		{-30, -180, 1, 0, 1},
		{0, -180, 1, 0, 1},
		{0, 0, 1, 1, 1},
	}

	for _, c := range cases {
		x, y := LatLong2XY(c.lat, c.long, c.zoom)
		if x != c.x {
			t.Errorf("LatLong2XY(%v, %v, %v) get(%v, %v), want(%v, %v)", c.lat, c.long, c.zoom, x, y, c.x, c.y)
		}
		if y != c.y {
			t.Errorf("LatLong2XY(%v, %v, %v) get(%v, %v), want(%v, %v)", c.lat, c.long, c.zoom, x, y, c.x, c.y)
		}
	}
}
