package main

import (
	"database/sql"
	"image"
	"log"
	"os"

	"github.com/billbell73/mosaic/lib/imageAnalyser"
)

const (
	tileSize = 20
	roughTilesAcross = 40
	hugeDistance = 10000000000
)

var tileColorAverages *map[string][3]int

// loads filenames and average colors of tiles from database
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

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

// iterates over original image and finds appropriate tile for each section
func createMosaic(original image.Image) ([][]string) {
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
	sectionSize := (bounds.Max.X - bounds.Min.X) / roughTilesAcross
	log.Println("sectionSize: ", sectionSize)

	for y := bounds.Min.Y; y + sectionSize <= bounds.Max.Y; y = y + sectionSize {
		var mosaicRow []string
		for x := bounds.Min.X; x + sectionSize <= bounds.Max.X; x = x + sectionSize {
			sectionCoords := image.Rect(x, y, x+sectionSize, y+sectionSize)
			section := original.(subImager).SubImage(sectionCoords)
			sectionColor := imageAnalyser.AverageRGB(section, imageAnalyser.TotalRGB)

			nearest := nearest(sectionColor, tileColorAverages)
			awsLinkPrefix := awsUrl + "/" + awsBucket + "/"
			mosaicRow = append(mosaicRow, awsLinkPrefix + nearest)
		}
		mosaic = append(mosaic, mosaicRow)
	}
	return mosaic
}

func width(mosaic [][]string) int {
	return len(mosaic[0]) * tileSize
}

// finds the tile that is the nearest match for a color
func nearest(target [3]int, colorAverages *map[string][3]int) string {
	var filename string
	smallest := hugeDistance

	for k, v := range *colorAverages {
		distance := distance(target, v)
		if distance < smallest {
			filename, smallest = k, distance
		}
	}
	// delete(*colorAverages, filename)
	return filename
}

// finds the squared Eucleadian distance between 2 points in 3D space
func distance(p1 [3]int, p2 [3]int) int {
	return sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2])
}

func sq(n int) int {
	return n * n
}
