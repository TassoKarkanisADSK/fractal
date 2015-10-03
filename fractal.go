package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"os"
)

const maxMagnitude2 = 1000
const maxIters = 2000

func iter(z, c complex128) complex128 {
	return z*z + c
}

func abs2(c complex128) float64 {
	a := real(c)
	b := imag(c)
	return a*a + b*b
}

func zero() complex128 {
	return 0
}

func compute(c complex128) int {
	z := zero()
	var iters int
	for {
		z = iter(z, c)
		iters++
		if abs2(z) > maxMagnitude2 {
			break
		}

		if iters == maxIters {
			break
		}
	}

	return iters
}

func getColor(a, total int64) color.Color {
	// just compute a red component for now
	r := float64(a)/float64(total)
	
	// index := int(255*r)
	// c := palette.Plan9[index]
	
	var c color.Color
	if a == total {
		c = color.Black
	} else {
		index := int(215*r)
		c = palette.WebSafe[index]
	}

	// var c color.RGBA
	// c.R = 0
	// c.G = uint8(255 * r)
	// c.B = uint8(255 * r)
	// c.A = 255

	return c
}

func computeImage(xres, yres int, x0, y0, step float64) image.Image {
	// allocate storage for the histogram
    histogram := make([]int, maxIters + 1)

	// allocate storage for the iteration counts
	counts := make([]int, xres*yres)

	// fill it out the histogram and iteration counts
	for j := 0; j < yres; j++ {
		y := y0 + float64(j)*step
		for i := 0; i < xres; i++ {
			x := x0 + float64(i)*step
			c := complex(x, y)
			iters := compute(c)
			histogram[iters]++
			counts[j*xres + i] = iters
		}
	}

	// compute the histogram total and partial sums
	partialSums := make([]int64, maxIters + 1)
	var total int64
	for i := 0; i < maxIters + 1; i++ {
		total += int64(histogram[i])
		partialSums[i] = total
	}

	// allocate the image
	rect := image.Rect(0, 0, xres, yres)
	img := image.NewRGBA(rect)
	
	// render the image based on the counts
	for j := 0; j < yres; j++ {
		for i := 0; i < xres; i++ {
			iters := counts[j*xres + i]
			col := getColor(partialSums[iters], total)
			img.Set(i, j, col)
		}
	}

	return img
}

func main() {
	var _ = fmt.Printf

	const XRes = 1024
	const YRes = 768

	const X = -0.65
	const Y = 0.475
	const Mag = 512

	step := 1.0 / (Mag * float64(XRes))
	x0 := float64(X) - 0.5*(float64(XRes)*step)
	y0 := float64(Y) - 0.5*(float64(YRes)*step)

	img := computeImage(XRes, YRes, x0, y0, step)
	w, _ := os.Create("fractal.png")
	png.Encode(w, img)
}
