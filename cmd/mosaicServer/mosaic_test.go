package main

import (
	"os"
	"image"
	"image/color"
	"testing"
)

func TestSq(t *testing.T) {
	var v int
	v = sq(4)
	if v != 16 {
		t.Error("Expected 16, got ", v)
	}
}

func TestDistance(t *testing.T) {
	var v int
	p1 := [3]int{1, 1, 2}
	p2 := [3]int{3, 4, 8}
	v = distance(p1, p2)
	if v != 49 {
		t.Error("Expected 49, got ", v)
	}
}

func TestNearest(t *testing.T) {
	target := [3]int{1, 1, 1}
	db := map[string][3]int{
		"image2": [3]int{10, 10, 10},
		"image1": [3]int{3, 2, 2},
		"image3": [3]int{3, 4, 19},
	}
	filename := nearest(target, &db)
	if filename != "image1" {
		t.Error("Expected \"image1\", got ", filename)
	}
}

//test helper method that colors a section of 'original' test image
func testPaint(img *image.NRGBA64, pt image.Point, sectionSize int, rgb [3]uint16) {
	for y := pt.Y; y < pt.Y + sectionSize ; y = y + 1 {
		for x := pt.X; x < pt.X + sectionSize; x = x + 1 {
			img.Set(x, y, color.NRGBA64{rgb[0], rgb[1], rgb[2], 60000})
		}
	}
}

func TestCreateMosaic(t *testing.T) {
	os.Setenv("AWS_URL", "https://url.com")
	os.Setenv("AWS_BUCKET", "bucket")

	tileColorAverages = &map[string][3]int{
		"image1.jpg": [3]int{3, 2, 2},
		"image2.jpg": [3]int{300, 300, 300},
		"image3.jpg": [3]int{10000, 10000, 10000},
	}

	r := image.Rectangle{image.Point{0, 0}, image.Point{250, 250}}
	original := image.NewNRGBA64(r)

	sectionSize := 250 / roughTilesAcross
	testPaint(original, image.Point{60, 60}, sectionSize, [3]uint16{50, 50, 50})
	testPaint(original, image.Point{120, 120}, sectionSize, [3]uint16{250, 350, 300})
	testPaint(original, image.Point{60, 0}, sectionSize, [3]uint16{9000, 9000, 9000})

	mosaic := createMosaic(original)

	tilesAcross := 250 / sectionSize
	if len(mosaic[0]) != tilesAcross {
		t.Error("Expected", tilesAcross, "tiles horizontally, got", len(mosaic[0]))
	}

	if mosaic[10][10] != "https://url.com/bucket/image1.jpg" {
		t.Error("Expected 'https://url.com/bucket/image1.jpg', got", mosaic[10][10])
	}

	if mosaic[20][20] != "https://url.com/bucket/image2.jpg" {
		t.Error("Expected 'https://url.com/bucket/image2.jpg', got", mosaic[20][20])
	}

	if mosaic[0][10] != "https://url.com/bucket/image3.jpg" {
		t.Error("Expected 'https://url.com/bucket/image3.jpg', got", mosaic[0][10])
	}
}

func TestWidth(t *testing.T) {
	mosaicRowLength := 3

	mosaic := make([][]string, 5)
	mosaic[0] = make([]string, mosaicRowLength)
	width := width(mosaic)

	if width !=  mosaicRowLength * tileSize {
		t.Error("Expected width", mosaicRowLength * tileSize, "got", width)
	}
}
