package main

import (
	"image"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/lib/pq"

	_ "image/jpeg"
)

func newHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/new.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10485760) // max body in memory is 10MB
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	original, format, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Format: ", format)
	log.Println("Bounds of original image: ", original.Bounds())

	mosaic := createMosaic(original)

	t, err := template.ParseFiles("views/show.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, mosaic)
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
