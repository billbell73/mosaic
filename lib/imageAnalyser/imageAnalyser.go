package imageAnalyser

import (
	"image"
	"image/color"
)

type mosaicImage interface {
	Bounds() image.Rectangle
	At(x, y int) color.Color
}

func AverageColor(img mosaicImage, fn colorTotaller) [3]int {
	bounds := img.Bounds()
	colors := fn(img, bounds)

	totalX := bounds.Max.X - bounds.Min.X
	totalY := bounds.Max.Y - bounds.Min.Y
	totalPixels := totalX * totalY

	return [3]int{
		colors[0] / totalPixels,
		colors[1] / totalPixels,
		colors[2] / totalPixels,
	}
}

type colorTotaller func(mosaicImage, image.Rectangle) [3]int

func TotalRGB(img mosaicImage, bounds image.Rectangle) [3]int {
	r, g, b := 0, 0, 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r, g, b = r+int(r1), g+int(g1), b+int(b1)
		}
	}
	return [3]int{r, g, b}
}
