package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func tileFilenames(dirPath string) []string {
	files, err := ioutil.ReadDir(dirPath)
	checkErr(err)

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	return filenames
}

type extractor func(string) (image.Image, string, error)

func openAndDecode(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(f)
}

type fileCreator func(string, image.Image) error

func writeImageToFile(filename string, img image.Image) error {
	out, err := os.Create(filename)
	checkErr(err)
	defer out.Close()

	return jpeg.Encode(out, img, nil)
}

func resizeImages(sourceDir string, destDir string, tileSize uint, filenames []string, extractFn extractor, createFn fileCreator) int {
	var resizedCount int
	for _, filename := range filenames {
		img, _, err := extractFn(sourceDir + filename)
		if err != nil {
			fmt.Printf("Couldn't decode file: \"%s\", Error: \"%s\"\n", filename, err)
			continue
		}
		tile := resize.Thumbnail(tileSize, tileSize, img, resize.Lanczos3)
		destPath := destDir + filename
		err = createFn(destPath, tile)
		checkErr(err)
		resizedCount++
	}
	return resizedCount
}

func main() {
	sourceDir := flag.String("source", "../natgeo_orig/", "relative path to source directory")
	destDir := flag.String("dest", "../natgeo_tiny/", "relative path to destination directory")
	tileSize := flag.Uint("pixels", 20, "dimensions of resized image in pixels")
	flag.Parse()

	filenames := tileFilenames(*sourceDir)

	resizedCount := resizeImages(*sourceDir, *destDir, *tileSize, filenames, openAndDecode, writeImageToFile)
	log.Printf("No. of images resized: %d", resizedCount)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
