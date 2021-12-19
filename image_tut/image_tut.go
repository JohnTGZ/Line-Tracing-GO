package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	width := 200
	height := 100

	topLeft := image.Point{0, 0}
	btmRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{topLeft, btmRight})

	cyan := color.RGBA{100, 200, 200, 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: //upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: //lower right quadrant
				img.Set(x, y, color.White)
			default:
				//use zero value
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)

}
