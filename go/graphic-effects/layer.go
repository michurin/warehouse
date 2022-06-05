package main

import (
	"image"
	"image/color"
)

// very naive impleventation!
// no bound checking
// no integrity checking

type Layers struct {
	W  int
	H  int
	R  [][]float64
	G  [][]float64
	B  [][]float64
	Y  [][]float64
	Cb [][]float64
	Cr [][]float64
}

func New(s image.Image) Layers {
	bn := s.Bounds()
	x0 := bn.Min.X
	x1 := bn.Max.X
	y0 := bn.Min.Y
	y1 := bn.Max.Y
	xm := x1 - x0
	ym := y1 - y0
	vr := make([][]float64, ym)
	vg := make([][]float64, ym)
	vb := make([][]float64, ym)
	cy := make([][]float64, ym)
	cb := make([][]float64, ym)
	cr := make([][]float64, ym)
	for y := 0; y < ym; y++ {
		vr[y] = make([]float64, xm)
		vg[y] = make([]float64, xm)
		vb[y] = make([]float64, xm)
		cy[y] = make([]float64, xm)
		cb[y] = make([]float64, xm)
		cr[y] = make([]float64, xm)
		for x := 0; x < xm; x++ {
			r, g, b, _ := s.At(x0+x, y0+y).RGBA() // [0..0xffff]
			vr[y][x] = float64(r) / 65536
			vg[y][x] = float64(g) / 65536
			vb[y][x] = float64(b) / 65536
			xcy, xcb, xcr := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			cy[y][x] = float64(xcy) / 256
			cb[y][x] = float64(xcb) / 256
			cr[y][x] = float64(xcr) / 256
		}
	}
	return Layers{
		W:  xm,
		H:  ym,
		R:  vr,
		G:  vg,
		B:  vb,
		Y:  cy,
		Cb: cb,
		Cr: cr,
	}
}
