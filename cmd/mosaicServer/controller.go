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

	_ "image/png"

	_ "github.com/lib/pq"
)

func newHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/new.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	r.ParseMultipartForm(10485760) // max body in memory is 10MB
	file, _, err := r.FormFile("image")
	log.Printf("%T", file)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	t05 := time.Now()

	original, format, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Format: ", format)
	log.Println("Bounds of original image: ", original.Bounds())

	t1 := time.Now()

	mosaic, width := createMosaic(original)

	t2 := time.Now()

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, original, nil); err != nil {
		log.Fatal("Unable to encode original image: ", err)
	}
	str := base64.StdEncoding.EncodeToString(buffer.Bytes())

	t3 := time.Now()

	log.Println("Parse form:", t05.Sub(t0))
	log.Println("decode file:", t1.Sub(t05))
	log.Println("Create mosaic:", t2.Sub(t1))
	log.Println("Encode original:", t3.Sub(t2))

	t, err := template.ParseFiles("views/show.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, show{mosaic, str, "2", width})
}

type show struct {
	Mosaic   [][]string
	Original string
	Duration string
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

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.ListenAndServe(":"+port, nil)
}
