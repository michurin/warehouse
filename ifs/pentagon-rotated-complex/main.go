package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"iter"
	"log"
	"math"
	"math/cmplx"
	"os"
)

var shifts = func() [5]complex128 {
	s := [5]complex128{}
	for i := range 5 {
		s[i] = cmplx.Rect(.6, math.Pi*float64(4*i-1)/10)
	}
	return s
}()

var rotate = cmplx.Rect(2.497, .384)

func all(v complex128, n int) iter.Seq[complex128] {
	return func(yield func(complex128) bool) {
		if n <= 0 {
			yield(v)
		} else {
			// finding two closest and iterate over them
			r0 := math.MaxFloat64
			r1 := math.MaxFloat64
			v0 := complex(0, 0)
			v1 := complex(0, 0)
			for _, s := range shifts {
				c := v + s
				r := cmplx.Abs(c)
				if r < r0 {
					r1 = r0
					r0 = r
					v1 = v0
					v0 = c
				} else if r < r1 {
					r1 = r
					v1 = c
				}
			}
			n--
			for u := range all(v0*rotate, n) {
				if !yield(u) {
					return
				}
			}
			for u := range all(v1*rotate, n) {
				if !yield(u) {
					return
				}
			}
		}
	}
}

func value(c complex128) float64 {
	// find closes
	r := math.MaxFloat64
	for u := range all(c, 6) { // 5 is depth
		q := cmplx.Abs(u)
		if q < r {
			r = q
		}
	}

	return r
}

const width, height = 768, 768

func colorNRGBA(v float64) color.NRGBA {
	switch {
	case v < 0:
		return color.NRGBA{0, 0, 0, 255}
	case v < 1:
		c := uint8(v * 256)
		return color.NRGBA{c, c, c, 255}
	case v < 2:
		return color.NRGBA{255, 255, 255, 255}
	case v < 3:
		c := uint8((2 - v) * 256)
		return color.NRGBA{c, 255, c, 255}
	case v < 4:
		c := uint8((v - 4) * 256)
		return color.NRGBA{c, 255, c, 255}
	default:
		return color.NRGBA{255, 255, 255, 255}
	}
	c := uint8(0)
	if v > 0 {
		v *= 255
		if v < 256 {
			c = uint8(v)
		} else {
			c = uint8(255)
		}
	}
	return color.NRGBA{
		R: c,
		G: c,
		B: c,
		A: 255,
	}
}

func main() {
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := range height {
		fmt.Printf("%6d\033[6D", y)
		for x := range width {
			img.Set(x, y, colorNRGBA(value(complex(float64(x*2-width)/width, float64(height-y*2)/height))))
		}
	}
	fmt.Println()

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
