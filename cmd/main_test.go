package main

import (
	"image"
	"image/color"
	"testing"
)

func TestSq(t *testing.T) {
	var v float64
	v = sq(4)
	if v != 16 {
		t.Error("Expected 16, got ", v)
	}
}

func TestDistance(t *testing.T) {
	var v float64
	p1 := [3]float64{1, 1, 2}
	p2 := [3]float64{3, 4, 8}
	v = distance(p1, p2)
	if v != 7 {
		t.Error("Expected 7, got ", v)
	}
}

type testMosaicImage struct{}
type testColor struct{}

func (t testMosaicImage) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{2, 2}}
}

func (t testColor) RGBA() (uint32, uint32, uint32, uint32) {
	return 3, 3, 3, 3
}

func (t testMosaicImage) At(x, y int) color.Color {
	return testColor{}
}

func stubTotalRGB(img mosaicImage, bounds image.Rectangle) [3]float64 {
	return [3]float64{36, 36, 36}
}

func TestAverageColor(t *testing.T) {
	v := averageColor(testMosaicImage{}, stubTotalRGB)
	if v != [3]float64{4, 4, 4} {
		t.Error("Expected [4 4 4], got ", v)
	}
}

func TestTotalRGB(t *testing.T) {
	mosaic := testMosaicImage{}
	bounds := mosaic.Bounds()
	v := totalRGB(mosaic, bounds)
	if v != [3]float64{27, 27, 27} {
		t.Error("Expected [27 27 27], got ", v)
	}
}

func TestNearest(t *testing.T) {
	target := [3]float64{1,1,1}
	db := map[string][3]float64{
		"image2": [3]float64{10,10,10},
		"image1": [3]float64{3,2,2},
		"image3": [3]float64{3,4,19},
	}
	filename := nearest(target, &db)
	if filename != "image1" {
		t.Error("Expected \"image1\", got ", filename)
	}
}

func stubOpenAndDecode(filePath string) (mosaicImage, string, error) {
	return testMosaicImage{}, "", nil
}

func TestTilesDB(t *testing.T) {
	filenames := []string{"a.jpg", "b.jpg", "c.jpg"}
	db := tilesDB("", filenames, stubOpenAndDecode)

 	for _, filename := range filenames {
     if db[filename] != [3]float64{3,3,3} {
     	t.Error("Expected blah, got ", db[filename])
     }
  }
}
