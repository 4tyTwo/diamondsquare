package main

import (
	"math"
	"math/rand"
)

// Generator generates random2d heightmap using square-diamond algorithm
type Generator struct {
	Roughness           float64
	Height, Width       int
	stepSize            int
	hmap                [][]float64
	currRoughness       float64
	genHeight, genWidth int
	transposed          bool
}

// Generate creates heightmap
func (gen *Generator) Generate() [][]float64 {
	gen.genHeight, gen.genWidth, gen.transposed = getGenerationSize(gen.Height, gen.Width)
	gen.hmap = createHeightMap(gen.genHeight, gen.genWidth)
	defer gen.eraseMap()
	for i := 0; i < iterationsCount(gen.genHeight); i++ {
		gen.currRoughness = math.Pow(gen.Roughness, float64(i))
		gen.stepSize = (gen.genWidth - 1) / int(math.Pow(2, float64(i)))
		gen.diamond()
		gen.square()
	}
	if gen.transposed {
		gen.hmap = transpose(gen.hmap)
	}
	return truncateHeightMap(gen.hmap, gen.Height, gen.Width)
}

func (gen *Generator) eraseMap() {
	gen.hmap = nil
}

func (gen *Generator) square() {
	halfStep, quaterStep := gen.stepSize/2, gen.stepSize/4
	for i := 0; i < gen.genHeight; i += halfStep {
		for j := quaterStep; j < gen.genWidth; j += halfStep {
			gen.hmap[i][j] = gen.squareDiscplace(i, j)
		}
	}
	for i := quaterStep; i < gen.genHeight; i += halfStep {
		for j := 0; j < gen.genWidth; j += halfStep {
			gen.hmap[i][j] = gen.squareDiscplace(i, j)
		}
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
	if i+quaterStep < gen.genHeight {
		total += gen.hmap[i+quaterStep][j]
		count++
	}
	if j+quaterStep < gen.genWidth {
		total += gen.hmap[i][j+quaterStep]
		count++
	}
	avg := total / float64(count)
	return weightedAverage(gen.currRoughness, rand.Float64(), avg)
}

func (gen *Generator) diamond() {
	halfStep, quaterStep := gen.stepSize/2, gen.stepSize/4
	for i := quaterStep; i < gen.genHeight; i += halfStep {
		for j := quaterStep; j < gen.genWidth-1; j += halfStep {
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
}

func weightedAverage(val1Weight, val1, val2 float64) float64 {
	return val1Weight*val1 + val2*(1-val1Weight)
}

func getGenerationSize(height, width int) (int, int, bool) {
	// biggest power of n
	var transposed bool
	if height > width {
		transposed = true
		width, height = height, width
	}
	logWidth := math.Ceil(math.Log2(float64(width - 1)))
	logHeight := math.Ceil(math.Log2(float64(height - 1)))
	if logWidth <= logHeight {
		logWidth++
	}
	return int(math.Pow(2, float64(logHeight))) + 1, int(math.Pow(2, float64(logWidth))) + 1, transposed
}

func transpose(m [][]float64) [][]float64 {
	// https://rosettacode.org/wiki/Matrix_transposition#2D_representation
	r := make([][]float64, len(m[0]))
	for x := range r {
		r[x] = make([]float64, len(m))
	}
	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}

func truncateHeightMap(hmap [][]float64, height, width int) [][]float64 {
	truncated := hmap[:height]
	for i := range truncated {
		truncated[i] = hmap[i][:width]
	}
	return truncated
}
