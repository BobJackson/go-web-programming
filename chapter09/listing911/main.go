package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

const currentPath = "chapter09/listing911"

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(currentPath + "/public"))
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

func upload(w http.ResponseWriter, _ *http.Request) {
	// print current path
	//_, _ = fmt.Fprintf(w, "Current path: %s\n", os.Getenv("PWD"))

	t, err := template.ParseFiles(currentPath + "/upload.html")
	if err != nil {
		fmt.Println(err)
		return
	}
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
			file, err := os.Open(currentPath + "/" + nearestFile)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil {
					t := resize(img, tileSize)
					tile := t.SubImage(t.Bounds())
					tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
					draw.Draw(newImage, tileBounds, tile, sp, draw.Src)
				} else {
					fmt.Println("Decode error:", err, nearestFile)
				}
			} else {
				fmt.Println("Open error:", nearestFile)
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
	mosaicStr := base64.StdEncoding.EncodeToString(buf2.Bytes())
	t1 := time.Now()
	images := map[string]string{
		"original": originalStr,
		"mosaic":   mosaicStr,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	t, _ := template.ParseFiles(currentPath + "/results.html")
	_ = t.Execute(w, images)
}
