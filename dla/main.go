package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

const (
	StickinessRadius  = float64(1)
	StickinessRadius2 = StickinessRadius * StickinessRadius
)

type (
	point    [3]float64
	pointKey [3]int
	label    struct {
		p   point
		age int
	}
)

func pointDiff(a, b point) point {
	return point{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func pointRadius2(p point) float64 {
	return p[0]*p[0] + p[1]*p[1] + p[2]*p[2]
}

func pointMul(p point, k float64) point {
	return point{p[0] * k, p[1] * k, p[2] * k}
}

func normolize(p point) point {
	r := math.Sqrt(pointRadius2(p))
	return point{p[0] / r, p[1] / r, p[2] / r}
}

func randomUnionPoint() point {
	for {
		p := point{rand.Float64()*2 - 1, rand.Float64()*2 - 1, rand.Float64()*2 - 1}
		if pointRadius2(p) < 1 {
			return p
		}
	}
}

type area struct {
	a map[pointKey][]label
}

func (a *area) key(p point) pointKey {
	return pointKey{
		int(math.Floor(p[0] / StickinessRadius)),
		int(math.Floor(p[1] / StickinessRadius)),
		int(math.Floor(p[2] / StickinessRadius)),
	}
}

func (a *area) set(p point, age int) {
	k := a.key(p)
	a.a[k] = append(a.a[k], label{p: p, age: age})
}

func (a *area) near(p point) bool {
	k := a.key(p)
	for dx := k[0] - 1; dx <= k[0]+1; dx++ {
		for dy := k[1] - 1; dy <= k[1]+1; dy++ {
			for dz := k[2] - 1; dz <= k[2]+1; dz++ {
				for _, e := range a.a[pointKey{dx, dy, dz}] {
					if pointRadius2(pointDiff(p, e.p)) < StickinessRadius2 {
						return true
					}
				}
			}
		}
	}
	return false
}

func (a *area) dump() []label {
	r := []label(nil)
	for _, v := range a.a {
		r = append(r, v...)
	}
	return r
}

func listOfPoints(n int) []label {
	defer func() {
		fmt.Println()
	}()
	startR := StickinessRadius
	ar := area{a: map[pointKey][]label{}}
	ar.set(point{}, 0)
	for i := 0; i < n; i++ {
		startR2 := startR * startR
		boundR2 := 4 * startR2
		p := pointMul(normolize(randomUnionPoint()), startR)
		for {
			r := pointRadius2(p)
			if r > boundR2 {
				p = pointMul(normolize(randomUnionPoint()), startR) // reset
			}
			if r < startR2 {
				if ar.near(p) {
					ar.set(p, i)
					startR = math.Max(startR, math.Sqrt(pointRadius2(p))+StickinessRadius)
					fmt.Print(i, n, i*100/n, "% R=", startR, p, "\033[K\033[G")
					break
				}
			}
			p = pointDiff(p, randomUnionPoint()) // it is add operation
		}
	}
	return ar.dump()
}

func colorr(x, a int, za, z, zb float64) color.RGBA { // TODO rename function
	k := (z - za) / (zb - za)
	t := float64(511*x) / float64(a)
	if t < 256 {
		return color.RGBA{
			R: 0,
			G: uint8(t * k),
			B: 0,
			A: 255,
		}
	}
	t -= 256
	return color.RGBA{
		R: uint8(t * k),
		G: uint8(t * k),
		B: 0,
		A: 255,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	amp := 1000000
	pp := listOfPoints(amp)
	sort.Slice(pp, func(i, j int) bool { return pp[i].p[2] < pp[j].p[2] })
	fmt.Println(pp)
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: -240,
			Y: -240,
		},
		Max: image.Point{
			X: 240,
			Y: 240,
		},
	})
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	za := pp[0].p[2]
	zb := pp[len(pp)-1].p[2]
	for _, p := range pp {
		c := colorr(p.age, amp, za, p.p[2], zb)
		img.Set(int(p.p[0]), int(p.p[1]), c)
	}
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
