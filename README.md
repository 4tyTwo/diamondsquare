# diamondsquare
## Description
This package implements Diamond-square algorithm for 2D noise generation in golang.  
More about the algorithm at [wikipedia](https://en.wikipedia.org/wiki/Diamond-square_algorithm)
## Usage
### Generator struct
Generator describes input parameters, those are: 
* Height - height (number of rows) of generated noise matrix
* Width  - width (number of collums) of generated noise matrix
* Roughness - value of "randomness", the value must be between 0 and 1. Greater values correspond to greater weight of randomness on each iteration of the algorithm.
### Generating noise
In order to generate actual noise you should either create Generator struct and call Generate method of it
```
gen := Generator{Roughness: roughness, Width: width, Height: height}
noise := gen.Generate()
```
or call Generate function with desired parameters of Generator
```
noise := diamondsquare.Generate(width, height, roughness)
```
### Caveats
Consider that this implementation works only with `(2^(n+1) + 1) x (2^n + 1)` matrices, so if you call `Generate` with different height/width ratio, you actually generate bigger matrix, but get it croped, so using generator to create matrices with high sides ratio will lead to significant generation time and excessive memory usage.
