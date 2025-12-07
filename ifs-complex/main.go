package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

func log(a ...any) {
	if !true {
		fmt.Println(a...)
	}
}

func step(n int, scale float64, trans [][2]complex128, c complex128) float64 {
	// TODO: caching
	log(n, "enter", c)
	if n <= 0 { // TODO: manage scale
		log(n, "return abs:", cmplx.Abs(c))
		return cmplx.Abs(c)
	} else {
		log(n, "loop", len(trans))
		r := math.MaxFloat64
		n--
		// TODO: consider overlapping:
		// - apply transformation
		// - find two(!) closest
		// - preform `step` to both two
		// - find the closest result of two `step` applications
		for _, t := range trans {
			q := step(n, scale*cmplx.Abs(t[1]), trans, (t[0]+c)*t[1])
			if q < r {
				r = q
			}
			log(n, "  q,r:", q, r)
		}
		return r
	}
}

func colorNRGBA(v float64) color.NRGBA {
	if v > 1.5 {
		return color.NRGBA{255, 255, 255, 255}
	}
	x := uint8((1 - math.Cos(v*2*math.Pi)) * 127.5)
	return color.NRGBA{x, x, x, 255}
}

func noerr(err error) {
	if err != nil {
		fmt.Println("ERROR", err.Error())
		os.Exit(2)
	}
}

func main() {

	if !true { // debug
		fmt.Println("Debug:", step(1, 1, [][2]complex128{
			{complex(0.5, -0.25), 1 / complex(.5, .25)},
			{complex(-0.5, -0.25), 1 / complex(-.5, .25)},
		}, complex(-.5, .25)))
		return
	}

	const iw = 512
	const ih = 512

	img := image.NewNRGBA(image.Rect(0, 0, iw, ih))

	trans := [][2]complex128{ // pairs: {shift, scale+rotate}
		// v1: equal scales, no overlapping: works perfect
		// {complex(0.5, -0.25), 1 / complex(.5, .25)},
		// {complex(-0.5, -0.25), 1 / complex(-.5, .25)},
		// v2: small inequality: scaling based artifacts [see TODO]
		// {complex(0.6, -0.3), 1 / complex(.4, .3)},
		// {complex(-0.4, -0.3), 1 / complex(-.6, .3)},
		// v3: going to be overlapped: artifacts show up [see TODO]
		{complex(0.6, -0.4), 1 / complex(.4, .4)},
		{complex(-0.4, -0.4), 1 / complex(-.6, .4)},
	}

	fmt.Println("0...")
	for y := range ih {
		for x := range iw {
			c := complex(float64(x*2-iw)/iw, float64(ih-y*2)/ih) * 2 // [-2, 2]
			img.Set(x, y, colorNRGBA(step(6, 1, trans, c)))
		}
		fmt.Printf("\033[1A\033[J%5.01f%% (height=%d)\n", (float64(y+1)*100)/ih, y+1)
	}

	f, err := os.Create("image.png")
	noerr(err)

	defer f.Close()

	noerr(png.Encode(f, img))

	fmt.Println("FIN.")
}
