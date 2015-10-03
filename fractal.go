package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

const maxMagnitude2 = 1000
const maxIters = 800

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

func getColor(iters int) color.RGBA {
	// just compute a red component for now
	max := float64(maxIters)
	it := float64(iters)
	r := (max - it) / max

	var c color.RGBA
	c.R = 0
	c.G = uint8(255 * r)
	c.B = uint8(255 * r)
	c.A = 255
	return c
}

func computeImage(xres, yres int, x0, y0, step float64) image.Image {
	// allocate the image
	rect := image.Rect(0, 0, xres, yres)
	img := image.NewRGBA(rect)

	// fill it out
	for j := 0; j < yres; j++ {
		y := y0 + float64(j)*step
		for i := 0; i < xres; i++ {
			x := x0 + float64(i)*step
			c := complex(x, y)
			iters := compute(c)
			col := getColor(iters)
			img.Set(i, j, col)
		}
	}

	return img
}

func main() {
	var _ = fmt.Printf

	const XRes = 1024
	const YRes = 768

	const X = 0.32
	const Y = 0.5
	const Mag = 20

	step := 1.0 / (Mag * float64(XRes))
	x0 := float64(X) - 0.5*(float64(XRes)*step)
	y0 := float64(Y) - 0.5*(float64(YRes)*step)

	img := computeImage(XRes, YRes, x0, y0, step)
	w, _ := os.Create("fractal.png")
	png.Encode(w, img)
}
