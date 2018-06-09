// Mercator projection algorithm get from https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Go
package main

import (
	"math"
)

func LatLong2XY(lat, long float64, z int) (x int, y int) {
	x = int(math.Floor((long + 180.0) / 360.0 * (math.Exp2(float64(z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(z)))))
	return
}

//func XY2LongLat(Long, Lat, z int) (Long float64, Lat float64) {
//	n := math.Pi - 2.0*math.Pi*float64(Lat)/math.Exp2(float64(z))
//	Lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
//	Long = float64(Long)/math.Exp2(float64(z))*360.0 - 180.0
//	return Lat, Long
//}
