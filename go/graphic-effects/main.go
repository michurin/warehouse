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
	sourceBounds := source.Bounds()
	fmt.Printf("Image: [%d-%d]x[%d-%d]\n", sourceBounds.Min.X, sourceBounds.Max.X, sourceBounds.Min.Y, sourceBounds.Max.Y)
	target := image.NewRGBA(image.Rect(sourceBounds.Min.X, sourceBounds.Min.Y, sourceBounds.Max.X-8, sourceBounds.Max.Y-8))
	w := [8]float64{-0.0625, -0.0625, -0.125, -0.25, 0.25, 0.125, 0.0625, 0.0625}
	for y := sourceBounds.Min.Y; y < sourceBounds.Max.Y-8; y++ {
		for x := sourceBounds.Min.X; x < sourceBounds.Max.X-8; x++ {
			sx := float64(0)
			sy := float64(0)
			for t := 0; t < 8; t++ {
				r, g, b, _ := source.At(x+t, y+3).RGBA()
				cy, _, _ := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
				sx += float64(cy) * w[t]
				r, g, b, _ = source.At(x+3, y+t).RGBA()
				cy, _, _ = color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
				sy += float64(cy) * w[t]
			}
			s := math.Hypot(sx, sy)
			s *= 2
			if s > 255 {
				s = 255
			}
			target.Set(x, y, color.RGBA{
				R: 0,
				G: uint8(s),
				B: 0,
				A: 0xff,
			})
		}
	}
	fh, err := os.Create("outimage-a-sobel.png")
	xerr(err)
	err = png.Encode(fh, target)
	xerr(err)
}
