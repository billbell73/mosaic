package main

import (
	"database/sql"
	"image"
	"log"
	"os"

	"github.com/billbell73/mosaic/lib/imageAnalyser"
)

const numberOfTilesRow = 40

var tileColorAverages *map[string][3]int

func init() {
	tileColorAverages = colorAverages()
}

func colorAverages() *map[string][3]int {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var (
		name  string
		red   int
		green int
		blue  int
	)

	colorAverages := make(map[string][3]int)

	rows, err := db.Query("select name, red, green, blue from averages")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &red, &green, &blue)
		if err != nil {
			log.Fatal(err)
		}
		colorAverages[name] = [3]int{red, green, blue}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return &colorAverages
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func createMosaic(original image.Image) [][]string {
	awsUrl := os.Getenv("AWS_URL")
	if awsUrl == "" {
		log.Fatal("$AWS_URL must be set")
	}
	awsBucket := os.Getenv("AWS_BUCKET")
	if awsBucket == "" {
		log.Fatal("$AWS_BUCKET must be set")
	}

	var mosaic [][]string
	bounds := original.Bounds()
	tileSize := bounds.Max.X / numberOfTilesRow
	log.Println("tileSize: ", tileSize)

	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + tileSize {
		var mosaicRow []string
		for x := bounds.Min.X; x < bounds.Max.X; x = x + tileSize {
			sectionCoords := image.Rect(x, y, x+tileSize, y+tileSize)
			section := original.(SubImager).SubImage(sectionCoords)
			sectionColor := imageAnalyser.AverageColor(section, imageAnalyser.TotalRGB)

			nearest := nearest(sectionColor, tileColorAverages)
			awsLinkPrefix := awsUrl + "/" + awsBucket + "/"
			mosaicRow = append(mosaicRow, awsLinkPrefix + nearest)
		}
		mosaic = append(mosaic, mosaicRow)
	}
	return mosaic
}

// find the nearest matching image
func nearest(target [3]int, db *map[string][3]int) string {
	var filename string
	var smallest int
	firstTime := true

	for k, v := range *db {
		distance := distance(target, v)
		if firstTime == true {
			filename, smallest = k, distance
			firstTime = false
		} else if distance < smallest {
			filename, smallest = k, distance
		}
	}
	// delete(*db, filename)
	return filename
}

// find the squared Eucleadian distance between 2 points
func distance(p1 [3]int, p2 [3]int) int {
	return sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2])
}

// find the square
func sq(n int) int {
	return n * n
}
