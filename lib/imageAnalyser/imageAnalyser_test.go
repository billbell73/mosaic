package imageAnalyser

import (
	"image"
	"image/color"
	"testing"
)

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

func stubTotalRGB(img mosaicImage, bounds image.Rectangle) [3]int {
	return [3]int{36, 36, 36}
}

func TestAverageRGB(t *testing.T) {
	v := AverageRGB(testMosaicImage{}, stubTotalRGB)
	if v != [3]int{9, 9, 9} {
		t.Error("Expected [9 9 9], got ", v)
	}
}

func TestTotalRGB(t *testing.T) {
	mosaic := testMosaicImage{}
	bounds := mosaic.Bounds()
	v := TotalRGB(mosaic, bounds)
	if v != [3]int{12, 12, 12} {
		t.Error("Expected [12 12 12], got ", v)
	}
}
