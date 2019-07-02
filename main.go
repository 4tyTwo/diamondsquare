package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func saveToFile(img *image.RGBA, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic("Can't create file")
	}
	png.Encode(f, img)
}

func createImage(buffer [][]float64) *image.RGBA {
	width := len(buffer[0])
	height := len(buffer)
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	var col color.RGBA
	for j := 0; j < width; j++ {
		for i := 0; i < height; i++ {
			baseColor := uint8(buffer[i][j] * 255)
			col = color.RGBA{baseColor, baseColor, baseColor, 255}
			img.Set(j, i, col)
		}
	}
	return img
}

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		os.Stderr.WriteString("Not enought arguments\nUsage: SquareDiamond width height path\n")
		panic("")
	}
	width, err := strconv.Atoi(args[0])
	if err != nil {
		os.Stderr.WriteString("Error parsing width")
		panic("")
	}
	height, err := strconv.Atoi(args[1])
	if err != nil {
		os.Stderr.WriteString("Error parsing height")
		panic("")
	}
	rand.Seed(time.Now().UnixNano())
	gen := Generator{Roughness: 0.45, Width: width, Height: height}
	hmap := gen.Generate()
	pathToFile := args[2]
	img := createImage(hmap)
	saveToFile(img, pathToFile)
}
