package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func averageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}
	totalPixels := float64(bounds.Max.X * bounds.Max.Y)
	return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}

func resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	ratio := bounds.Dx() / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.X/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)})
		}
	}
	return *out
}

func tilesDB() map[string][3]float64 {
	fmt.Println("Start populating tiles db ...")
	db := make(map[string][3]float64)
	files, _ := os.ReadDir("tiles")
	for _, f := range files {
		name := "tiles/" + f.Name()
		file, err := os.Open(name)
		if err == nil {
			img, _, err := image.Decode(file)
			if err == nil {
				db[name] = averageColor(img)
			} else {
				fmt.Println("error in populating TILEDB:", err, name)
			}
		} else {
			fmt.Println("cannot open file", name, err)
		}
		_ = file.Close()
	}
	fmt.Println("Finished populating tiles db.")
	return db
}

func nearest(target [3]float64, db *map[string][3]float64) string {
	var filename string
	smallest := 100_0000.0
	for k, v := range *db {
		dist := distance(target, v)
		if dist < smallest {
			filename, smallest = k, dist
		}
	}
	delete(*db, filename)
	return filename
}

func distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

func sq(x float64) float64 { return x * x }

var TILESDB map[string][3]float64

func cloneTilesDB() map[string][3]float64 {
	db := make(map[string][3]float64)
	for k, v := range TILESDB {
		db[k] = v
	}
	return db
}

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", upload)
	mux.HandleFunc("/mosaic", mosaic)
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	TILESDB = tilesDB()
	fmt.Println("Mosaic server started.")
	_ = server.ListenAndServe()
}

func upload(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("upload.html")
	_ = t.Execute(w, nil)
}

func mosaic(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()

	_ = r.ParseMultipartForm(10485760)
	file, _, _ := r.FormFile("image")
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))

	original, _, _ := image.Decode(file)
	bounds := original.Bounds()

	newImage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))

	db := cloneTilesDB()

	sp := image.Point{X: 0, Y: 0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += tileSize {
			r, g, b, _ := original.At(x, y).RGBA()
			rgb := [3]float64{float64(r), float64(g), float64(b)}

			nearestFile := nearest(rgb, &db)
			file, err := os.Open(nearestFile)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil {
					t := resize(img, tileSize)
					tile := t.SubImage(t.Bounds())
					tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
					draw.Draw(newImage, tileBounds, tile, sp, draw.Src)
				} else {
					fmt.Println("error:", err, nearestFile)
				}
			} else {
				fmt.Println("error:", nearestFile)
			}
			if file != nil {
				_ = file.Close()
			}
		}
	}
	buf1 := new(bytes.Buffer)
	_ = jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	buf2 := new(bytes.Buffer)
	_ = jpeg.Encode(buf2, newImage, nil)
	mosaic := base64.StdEncoding.EncodeToString(buf2.Bytes())
	t1 := time.Now()
	images := map[string]string{
		"original": originalStr,
		"mosaic":   mosaic,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	t, _ := template.ParseFiles("result.html")
	_ = t.Execute(w, images)
}