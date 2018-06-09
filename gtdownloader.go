package main

import (
	"flag"
	"fmt"
	"strings"
	"strconv"
	"os"
	"net/http"
	"image"
	"image/draw"
	"image/jpeg"
)

type coord struct {
	Lat  float64
	Long float64
}

type coords []coord

func (c *coords) String() string {
	var r string
	for i, v := range *c {
		if i != 0 {
			r += " "
		}
		r += fmt.Sprintf("%f,%f", v.Lat, v.Long)
	}

	return r
}

func (c *coords) Set(value string) error {
	arr := strings.Split(value, ",")
	lat, err := strconv.ParseFloat(arr[0], 64)
	exitIfError(err == nil, 1, "invalid marker location", err)
	exitIfError(-85.05113 <= lat && lat <= 85.05113, 1, "latitude must be between -85.05113 and 85.05113, but get", lat)

	long, err := strconv.ParseFloat(arr[1], 64)
	exitIfError(err == nil, 1, "invalid marker location", err)
	exitIfError(-180 <= long && long <= 180, 1, "longitude must be between -180 and 180, but get", lat)

	*c = append(*c, coord{Lat: lat, Long: long})
	return nil
}

func exitIfError(ok bool, code int, prompt string, msg interface{}) {
	if !ok {
		fmt.Fprintf(os.Stderr, "%s: %v\n", prompt, msg)
		os.Exit(code)
	}
}

func ensureCacheDir() {
	if _, err := os.Stat(".tiles"); os.IsNotExist(err) {
		os.Mkdir("tiles", os.ModeDir|0755)
	}
}

func parseArguments() (topLeft, bottomRight coord, zoom int, lyrs string, out string) {
	var markers coords
	flag.Var(&markers, "m", "markers coordinates, like 23.1,90.2")
	flag.IntVar(&zoom, "z", 0, "zoom level, between 0 and 21 inclusive")
	flag.StringVar(&lyrs, "lyrs", "y", "lyrs parameter, y=hybrid s=satellite t=train m=map")
	flag.StringVar(&out, "out", "out.jpg", "output filename")
	flag.Parse()

	exitIfError(0 <= zoom && zoom <= 21, 1, "zoom must be between 0 and 21, but get", zoom)

	minLong, maxLong, minLat, maxLat := 180.0, -180.0, 90.0, -90.0
	for _, c := range markers {
		if c.Long < minLong {
			minLong = c.Long
		}
		if c.Long > maxLong {
			maxLong = c.Long
		}
		if c.Lat < minLat {
			minLat = c.Lat
		}
		if c.Lat > maxLat {
			maxLat = c.Lat
		}
	}

	topLeft = coord{Lat: maxLat, Long: minLong}
	bottomRight = coord{Lat: minLat, Long: maxLong}
	return
}


func saveTile(x, y, zoom int, lyrs string) image.Image {
	// lyrs parameter are: y=hybrid s=satellite t=train m=map
	url := fmt.Sprintf("http://mt1.google.com/vt/lyrs=%s&x=%d&y=%d&z=%d", lyrs, x, y, zoom)
	response, err := http.Get(url)
	exitIfError(err == nil, 3, "fetch tile fail", err)

	defer response.Body.Close()

	tile, _, err := image.Decode(response.Body)
	exitIfError(err == nil, 3, "fetch tile fail", err)

	return tile
}

func main() {
	topLeft, bottomRight, zoom, lyrs, out := parseArguments()
	ensureCacheDir()

	x0, y0 := LatLong2XY(topLeft.Lat, topLeft.Long, zoom)
	x1, y1 := LatLong2XY(bottomRight.Lat, bottomRight.Long, zoom)

	tileCount := (x1 - x0 + 1) * (y1 - y0 + 1)
	exitIfError(tileCount < 1000, 1, "too many tiles, please decrease zoom or reduce the scope of markers", zoom)

	fmt.Println("fetching tiles...")
	const tileWidth, tileHeight = 256, 256
	w := (x1 - x0 + 1) * tileWidth
	h := (y1 - y0 + 1) * tileHeight
	dest := image.NewRGBA(image.Rect(0, 0, w, h))
	for c, x := 0, x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			tile := saveTile(x, y, zoom, lyrs)

			dx, dy := x-x0, y-y0
			dp := image.Point{X: dx * tileWidth, Y: dy * tileHeight}
			r := image.Rectangle{Min: dp, Max: dp.Add(tile.Bounds().Size())}
			draw.Draw(dest, r, tile, tile.Bounds().Min, draw.Src)

			c++
			fmt.Printf("%d/%d\n", c, tileCount)
		}
	}


	outfile, err := os.Create(out)
	exitIfError(err == nil, 3, "create output file fail", err)

	opt := jpeg.Options{Quality: 100}
	jpeg.Encode(outfile, dest, &opt)
}
