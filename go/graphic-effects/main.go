package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

func xerr(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func filename() string {
	if len(os.Args) != 2 {
		xerr(errors.New("It needs one argument: source image filename"))
	}
	return os.Args[1]
}

func main() {
	reader, err := os.Open(filename())
	xerr(err)
	source, _, err := image.Decode(reader)
	xerr(err)

	l := New(source)
	fmt.Printf("Image: %dx%d\n", l.R.Width(), l.R.Height())
	target := image.NewRGBA(image.Rect(0, 0, l.R.Width()-8, l.R.Height()-8))
	w := [8]float64{-0.0625, -0.0625, -0.125, -0.25, 0.25, 0.125, 0.0625, 0.0625}
	for y := 0; y < l.R.Height()-8; y++ {
		for x := 0; x < l.R.Width()-8; x++ {
			sx := float64(0)
			sy := float64(0)
			for t := 0; t < 8; t++ {
				sx += l.Y.At(x+t, y+3) * w[t]
				sy += l.Y.At(x+3, y+t) * w[t]
			}
			s := math.Hypot(sx, sy)
			s = .1 + .9*s
			target.Set(x, y, color.RGBA{
				R: uint8(s * l.R.At(x+3, y+3) * 256),
				G: uint8(s * 256),
				B: uint8(s * l.B.At(x+3, y+3) * 256),
				A: 0xff,
			})
		}
	}
	fh, err := os.Create("outimage-a-sobel.png")
	xerr(err)
	err = png.Encode(fh, target)
	xerr(err)
}
