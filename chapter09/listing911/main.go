package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
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
	log.Println("Mosaic server started.")
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

	_ = r.ParseMultipartForm(10485760) // max body in memory is 10MB
	file, _, _ := r.FormFile("image")
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))

	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	db := cloneTilesDB()

	c1 := cut(original, &db, tileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	c2 := cut(original, &db, tileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	c3 := cut(original, &db, tileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	c4 := cut(original, &db, tileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)
	c := combine(bounds, c1, c2, c3, c4)

	buf1 := new(bytes.Buffer)
	_ = jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	t1 := time.Now()
	images := map[string]string{
		"original": originalStr,
		"mosaic":   <-c,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	t, _ := template.ParseFiles(currentPath + "/results.html")
	_ = t.Execute(w, images)
}

func cut(original image.Image, db *DB, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
	c := make(chan image.Image)
	sp := image.Point{X: 0, Y: 0}
	go func() {
		newImage := image.NewNRGBA(image.Rect(x1, y1, x2, y2))
		for y := y1; y < y2; y += tileSize {
			for x := x1; x < x2; x += tileSize {
				r, g, b, _ := original.At(x, y).RGBA()
				color := [3]float64{float64(r), float64(g), float64(b)}
				nearest := db.nearest(color)
				file, err := os.Open(currentPath + "/" + nearest)
				if err != nil {
					log.Printf("Open error: %v", err)
					continue
				}
				img, _, err := image.Decode(file)
				_ = file.Close()
				if err != nil {
					log.Printf("Decode error: %v", err)
					continue
				}
				t := resize(img, tileSize)
				tile := t.SubImage(t.Bounds())
				tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
				draw.Draw(newImage, tileBounds, tile, sp, draw.Src)
			}
			c <- newImage.SubImage(newImage.Rect)
		}
	}()
	return c
}

func combine(r image.Rectangle, c1, c2, c3, c4 <-chan image.Image) <-chan string {
	c := make(chan string)
	go func() {
		var wg sync.WaitGroup
		img := image.NewNRGBA(r)
		copyImage := func(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
			draw.Draw(dst, r, src, sp, draw.Src)
			wg.Done()
		}
		wg.Add(4)
		var s1, s2, s3, s4 image.Image
		var ok1, ok2, ok3, ok4 bool
		for {
			select {
			case s1, ok1 = <-c1:
				go copyImage(img, s1.Bounds(), s1, image.Point{X: r.Min.X, Y: r.Min.Y})
			case s2, ok2 = <-c2:
				go copyImage(img, s2.Bounds(), s2, image.Point{X: r.Max.X / 2, Y: r.Min.Y})
			case s3, ok3 = <-c3:
				go copyImage(img, s3.Bounds(), s3, image.Point{X: r.Min.X, Y: r.Max.Y / 2})
			case s4, ok4 = <-c4:
				go copyImage(img, s4.Bounds(), s4, image.Point{X: r.Max.X / 2, Y: r.Max.Y / 2})
			}
			if ok1 && ok2 && ok3 && ok4 {
				break
			}
		}
		wg.Wait()
		buf2 := new(bytes.Buffer)
		_ = jpeg.Encode(buf2, img, nil)
		c <- base64.StdEncoding.EncodeToString(buf2.Bytes())
	}()
	return c
}
