package main

import (
	"fmt"
	"image"
	"os"
	"io/ioutil"
	// "bufio"
	// "image/draw"
	"image/color"
	"math"

	_ "image/jpeg"
)

type mosaicImage interface {
	Bounds() image.Rectangle
	At(x, y int) color.Color
}

// find the average color of the picture
func averageColor(img mosaicImage, fn colorTotaller) [3]float64 {
	bounds := img.Bounds()
	colors := fn(img, bounds)
	totalX := bounds.Max.X + 1
	totalY := bounds.Max.Y + 1
	totalPixels := float64(totalX * totalY)

	return [3]float64{
		colors[0] / totalPixels,
		colors[1] / totalPixels,
		colors[2] / totalPixels,
  }
}

type colorTotaller func(mosaicImage, image.Rectangle) [3]float64

func totalRGB(img mosaicImage, bounds image.Rectangle) [3]float64 {
	r, g, b := 0.0, 0.0, 0.0
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}
	return [3]float64{r, g, b}
}

type extractor func(string) (mosaicImage, string, error)

func openAndDecode(filename string) (mosaicImage, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(f)
}

// populate a tiles database in memory
func tilesDB(dirPath string, filenames []string, fn extractor) map[string][3]float64 {
    fmt.Println("Start populating tiles db ...")
    db := make(map[string][3]float64)

    for _, filename := range filenames {
      img, _, err := fn(dirPath + filename)
      if err == nil {
          db[filename] = averageColor(img, totalRGB)
      } else {
          fmt.Printf("Couldn't decode file: \"%s\", Error: \"%s\"\n", filename, err)
      }
    }

    fmt.Println("Finished populating tiles db.")
    return db
}

// find the nearest matching image
func nearest(target [3]float64, db *map[string][3]float64) string {
  var filename string
  var smallest float64
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

// type distancer func([3]float64, [3]float64) float64

// find the Eucleadian distance between 2 points
func distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

// find the square
func sq(n float64) float64 {
	return n * n
}

// resize an image to its new width
func resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	width := bounds.Max.X - bounds.Min.X
	ratio := width / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.X/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return *out
}

// func cloneTilesDB() map[string][3]float64 {
// 		original := tilesDB()
//     db2 := make(map[string][3]float64)
//     for k, v := range original {
//         db2[k] = v
//     }
//     return db2
// }

func main() {
	original, _, err := openAndDecode("../instagram_tiny/source.jpg")
	checkErr(err)
	// tileSize := 100
	bounds := original.Bounds()

	fmt.Println("bounds", bounds.Max.X)

	avg := averageColor(original, totalRGB)
	fmt.Println("avg", avg)


  dirPath := "../instagram_tiny/"
  files, err := ioutil.ReadDir(dirPath)

  var filenames []string
  for _, file := range files {
    filenames = append(filenames, file.Name())
  }


  db := tilesDB(dirPath, filenames, openAndDecode)
  fmt.Println("db: ", db)

	// create a new image for the mosaic
	//  newimage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))
	//  // build up the tiles database
	//  db := tilesDB()
	//  // source point for each tile, which starts with 0, 0 of each tile
	//  sp := image.Point{0, 0}
	//  fmt.Println("1")
	//  for y := bounds.Min.Y; y < bounds.Max.Y; y = y + tileSize {
	//    for x := bounds.Min.X; x < bounds.Max.X; x = x + tileSize {
	//    		fmt.Println("2")
	//        // use the top left most pixel as the average color
	//        r, g, b, _ := original.At(x, y).RGBA()
	//        color := [3]float64{float64(r), float64(g), float64(b)}
	//        // get the closest tile from the tiles DB
	//        nearest := nearest(color, &db)
	//        file, err := os.Open(nearest)
	//        checkErr(err)
	//        fmt.Println("3")
	//        if err == nil {
	//            img, _, err := image.Decode(file)
	//            if err == nil {
	//                // resize the tile to the correct size
	//                t := resize(img, tileSize)
	//                tile := t.SubImage(t.Bounds())
	//                tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
	//                // draw the tile into the mosaic
	//                draw.Draw(newimage, tileBounds, tile, sp, draw.Src)
	//                fmt.Println("4")
	//            } else {
	//                fmt.Println("error:", err, nearest)
	//            }
	//        } else {
	//            fmt.Println("error:", nearest)
	//        }
	//        file.Close()
	//    }
	//  }

	// fmt.Printf("hi %T\n", newimage)
	// // fmt.Println(newimage)

	// out, err := os.Create("../instagram_tiny/output.jpg")
	// checkErr(err)

	// var opt jpeg.Options

	//  opt.Quality = 80
	//  // ok, write out the data into the new JPEG file

	//  err = jpeg.Encode(out, newimage, &opt)
	//  checkErr(err)

	// func Encode(w io.Writer, m image.Image, o *Options)
	// checkErr(err)
	// fmt.Println(averageColor(image))

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
