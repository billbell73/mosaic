// creates photo-mosaics from uploaded photos
package main

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"time"
	"sync"
	"runtime"

	_ "image/png"

	_ "github.com/lib/pq"
)

const tenMB = 10485760

// renders form for uploading photo
func newHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/new.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

// renders view with photo-mosaic and original uploaded image
func showHandler(w http.ResponseWriter, r *http.Request) {
  runtime.GOMAXPROCS(2)
	t0 := time.Now()

	r.ParseMultipartForm(tenMB)
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	t1 := time.Now()

	original, format, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Format of original image: ", format)
	log.Println("Bounds of original image: ", original.Bounds())

	var wg sync.WaitGroup
	wg.Add(1)

	var originalEncoded string
	go func() {
		defer wg.Done()
		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, original, nil); err != nil {
			log.Fatal("Unable to encode original image: ", err)
		}
		originalEncoded = base64.StdEncoding.EncodeToString(buffer.Bytes())
	}()

	t2 := time.Now()

	mosaic := createMosaic(original)
	width := width(mosaic)

	t, err := template.ParseFiles("views/show.html")
	if err != nil {
		log.Fatal(err)
	}

	t3 := time.Now()
	wg.Wait()
	t4 := time.Now()

	t.Execute(w, show{mosaic, originalEncoded, width})

	log.Println("Parse form: ", t1.Sub(t0))
	log.Println("Create mosaic: ", t3.Sub(t2))
	log.Println("Encode original & create mosaic: ", t4.Sub(t2))
}

type show struct {
	Mosaic   [][]string
	Original string
	Width    int
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/new", http.StatusFound)
	})
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/show", showHandler)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	tileColorAverages = colorAverages()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.ListenAndServe(":"+port, nil)
}
