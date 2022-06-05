package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

// low level util

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

func black(img *image.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 255,
			})
		}
	}
}

func saveImage(img image.Image, filename string) {
	fh, err := os.Create(filename)
	xerr(err)
	err = png.Encode(fh, img)
	xerr(err)
}

// sobel (sort of, lightweight)

func sobel(v [][]float64, ow, oh int) ([][]float64, int, int) {
	wgh := [8]float64{-0.0625, -0.0625, -0.125, -0.25, 0.25, 0.125, 0.0625, 0.0625}
	cs := len(wgh)
	w := ow - cs
	h := oh - cs
	d := cs / 2
	r := make([][]float64, h)
	for y := 0; y < h; y++ {
		r[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			sx := float64(0)
			sy := float64(0)
			for t := 0; t < cs; t++ { // yes it is not full convolution, we consider just cross
				sx += v[y+d][x+t] * wgh[t]
				sy += v[y+t][x+d] * wgh[t]
			}
			r[y][x] = math.Hypot(sx, sy)
		}
	}
	return r, w, h
}

// random points

type Point struct{ X, Y float64 }

func randPoints(n int, v [][]float64, w, h int) []Point {
	type u struct {
		v float64
		p Point
	}
	a := make([]u, w*h)
	c := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			a[c] = u{v: math.Pow(v[y][x], .7), p: Point{X: float64(x), Y: float64(y)}}
			c++
		}
	}
	sort.Slice(a, func(i, j int) bool { return a[i].v < a[j].v })
	base := a[0].v
	mag := a[len(a)-1].v
	fmt.Printf("Magnitude: %f - %f\n", base, mag)
	p := make([]Point, n)
	for i := 0; i < n; i++ {
		rnd := rand.Float64()*(mag-base) + base
		idx := sort.Search(len(a), func(i int) bool { return rnd < a[i].v })
		// fmt.Println(rnd, idx, a[idx])
		p[i] = a[idx].p
		p[i].X += rand.Float64() - .5 // randomization
		p[i].Y += rand.Float64() - .5
	}
	return p
}

// demo images

func sobelImage(l Layers) {
	sl, sw, sh := sobel(l.Y, l.W, l.H)
	fmt.Printf("Image: %dx%d -> %dx%d\n", l.W, l.H, sw, sh)
	shift := (l.W - sw) / 2
	target := image.NewRGBA(image.Rect(0, 0, sw, sh))
	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			s := sl[y][x]
			m := .2 + .2*s
			v := s * 4
			if v > 1 {
				v = .999
			}
			target.Set(x, y, color.RGBA{
				R: uint8(m * l.R[y+shift][x+shift] * 256),
				G: uint8(v * 256),
				B: uint8(m * l.B[y+shift][x+shift] * 256),
				A: 0xff,
			})
		}
	}
	saveImage(target, "outimage-a-sobel.png")
}

func randPointsImage(l Layers) {
	const factor = 2
	sl, w, h := sobel(l.Y, l.W, l.H)
	target := image.NewRGBA(image.Rect(0, 0, w*factor, h*factor))
	black(target)
	pp := randPoints(10000, sl, w, h)
	for _, p := range pp {
		target.Set(int(p.X*factor), int(p.Y*factor), color.RGBA{
			R: 0,
			G: 255,
			B: 0,
			A: 255,
		})
	}
	saveImage(target, "outimage-b-random.png")
}

func voronoiImage(l Layers) {
	// very naive implementation; O(number of pivot points)
	const factor = 3
	const power = 4
	sl, w, h := sobel(l.Y, l.W, l.H)
	wt := w * factor
	ht := h * factor
	target := image.NewRGBA(image.Rect(0, 0, wt, ht))
	black(target)
	pp := randPoints(7000, sl, w, h)
	for y := 0; y < ht; y++ {
		fmt.Printf("Progress: %d/%d", y, ht)
		for x := 0; x < wt; x++ {
			nd := 1_000_000_000.0
			nx := 0
			ny := 0
			for _, p := range pp {
				dx := float64(x) - p.X*factor
				dy := float64(y) - p.Y*factor
				d := math.Pow((math.Pow(math.Abs(dx), power) + math.Pow(math.Abs(dy), power)), 1./power)
				if d < nd {
					nd = d
					nx = int(p.X)
					ny = int(p.Y)
				}
			}
			target.Set(x, y, color.RGBA{
				R: uint8(256 * l.R[ny][nx]),
				G: uint8(256 * l.G[ny][nx]),
				B: uint8(256 * l.B[ny][nx]),
				A: 255,
			})
		}
		fmt.Print("\x1b[0G\x1b[0K")
	}
	saveImage(target, "outimage-c-voronoi.png")
}

func main() {
	reader, err := os.Open(filename())
	xerr(err)
	source, _, err := image.Decode(reader)
	xerr(err)

	l := New(source)

	sobelImage(l)

	rand.Seed(time.Now().UnixNano())
	randPointsImage(l)

	voronoiImage(l)
}
