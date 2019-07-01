package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strconv"
)

// Generator generates random2d heightmap using square-diamond algorithm
type Generator struct {
	Roughness     float64
	Height, Width int
	stepSize      int
	hmap          [][]float64
	currRoughness float64
}

// Generate creates heightmap
func (gen *Generator) Generate() [][]float64 {
	gen.hmap = createHeightMap(gen.Height, gen.Width)
	defer gen.eraseMap()
	for i := 0; i < iterationsCount(gen.Height); i++ {
		gen.currRoughness = math.Pow(gen.Roughness, float64(i))
		gen.stepSize = (gen.Width - 1) / int(math.Pow(2, float64(i)))
		gen.diamond()
		gen.square()
	}
	return gen.hmap
}

func (gen *Generator) eraseMap() {
	gen.hmap = nil
}

func (gen *Generator) square() {
	halfStep, quaterStep := gen.stepSize/2, gen.stepSize/4
	for i := 0; i < gen.Height; i += halfStep {
		for j := quaterStep; j < gen.Width; j += halfStep {
			gen.hmap[i][j] = gen.squareDiscplace(i, j)
		}
	}
	for i := quaterStep; i < gen.Height; i += halfStep {
		for j := 0; j < gen.Width; j += halfStep {
			gen.hmap[i][j] = gen.squareDiscplace(i, j)
		}
	}
}

func printMap(hmap [][]float64) {
	for _, row := range hmap {
		for _, item := range row {
			fmt.Printf("%5.2f", item)
		}
		fmt.Println()
	}
}

func (gen *Generator) squareDiscplace(i, j int) float64 {
	total, count := 0.0, 0
	quaterStep := gen.stepSize / 4
	if i-quaterStep >= 0 {
		total += gen.hmap[i-quaterStep][j]
		count++
	}
	if j-quaterStep >= 0 {
		total += gen.hmap[i][j-quaterStep]
		count++
	}
	if i+quaterStep < gen.Height {
		total += gen.hmap[i+quaterStep][j]
		count++
	}
	if j+quaterStep < gen.Width {
		total += gen.hmap[i][j+quaterStep]
		count++
	}
	avg := total / float64(count)
	return weightedAverage(gen.currRoughness, rand.Float64(), avg)
}

func (gen *Generator) diamond() {
	halfStep, quaterStep := gen.stepSize/2, gen.stepSize/4
	for i := quaterStep; i < gen.Height; i += halfStep {
		for j := quaterStep; j < gen.Width-1; j += halfStep {
			gen.hmap[i][j] = gen.diamondDiscplace(i, j)
		}
	}
}

func (gen *Generator) diamondDiscplace(i, j int) float64 {
	quaterStep := gen.stepSize / 4
	ul := gen.hmap[i-quaterStep][j-quaterStep]
	ur := gen.hmap[i-quaterStep][j+quaterStep]
	ll := gen.hmap[i+quaterStep][j-quaterStep]
	lr := gen.hmap[i+quaterStep][j+quaterStep]
	avg := (ul + ur + ll + lr) / 4.0
	return weightedAverage(gen.currRoughness, rand.Float64(), avg)
}

func createHeightMap(height, width int) [][]float64 {
	hmap := make([][]float64, height)
	for i := range hmap {
		hmap[i] = make([]float64, width)
	}
	// initialize corners with initial values
	hmap[0][0] = rand.Float64()
	hmap[0][width-1] = rand.Float64()
	hmap[height-1][0] = rand.Float64()
	hmap[height-1][width-1] = rand.Float64()
	hmap[0][(width-1)/2] = rand.Float64()
	hmap[height-1][(width-1)/2] = rand.Float64()
	return hmap
}

func iterationsCount(height int) int {
	return int(math.Log2(float64(height - 1)))
	// TODO: check value here
}

func weightedAverage(val1Weight, val1, val2 float64) float64 {
	return val1Weight*val1 + val2*(1-val1Weight)
}

func saveToFile(img *image.RGBA, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic("Can't create file")
	}
	png.Encode(f, img)
}

func createImage(buffer [][]float64, width, height int) *image.RGBA {
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
	gen := Generator{Roughness: 0.5, Width: width, Height: height}
	hmap := gen.Generate()
	pathToFile := args[2]
	img := createImage(hmap, width, height)
	saveToFile(img, pathToFile)
}
