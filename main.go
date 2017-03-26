package main

import (
	"flag"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

var (
	rows       = flag.Int("led-rows", 32, "number of rows supported")
	parallel   = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain      = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness = flag.Int("brightness", 100, "brightness (0-100)")
)


type Circle struct {
	X, Y, R float64
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 127
	} else {
		return 255
	}
}


func main() {
	config := &rgbmatrix.DefaultConfig
	config.Rows = *rows
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)

	defer tk.Close()

	tk.PlayAnimation(NewAnimation(image.Point{64, 32}))
}


func init() {
	flag.Parse()
}



func fatal(err error) {
	if err != nil {
		panic(err)
	}
}


type Animation struct {
	ctx      *gg.Context
	position image.Point
	dir      image.Point
	stroke   int
}


func NewAnimation(sz image.Point) *Animation {
	return &Animation{
		ctx:    gg.NewContext(sz.X, sz.Y),
		dir:    image.Point{1, 1},
		stroke: 5,
	}
}


func (a *Animation) Next() (image.Image, <-chan time.Time, error) {


	a.ctx.SetColor(color.Black)
	a.ctx.Clear()


	w := a.ctx.Width()
	h := a.ctx.Height()



	var hw, hh float64 = float64(w / 2), float64(h / 2)
	circles := []*Circle{&Circle{}, &Circle{}, &Circle{}}
	steps := 20
	step :=0

	step += 1 
	step %= steps

	θ := 2.0 * math.Pi / float64(steps) * float64(step)
	for i, circle := range circles {
		θ0 := 2 * math.Pi / 3 * float64(i)
		circle.X = hw - 40*math.Sin(θ0) - 20*math.Sin(θ0+θ)
		circle.Y = hh - 40*math.Cos(θ0) - 20*math.Cos(θ0+θ)
		circle.R = 50
	}
	
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			a.ctx.SetColor(color.RGBA{
                                circles[0].Brightness(float64(x), float64(y)),
                                circles[1].Brightness(float64(x), float64(y)),
                                circles[2].Brightness(float64(x), float64(y)),
                                255,
                        })
			a.ctx.DrawCircle(float64(x),float64(y),1.0)
		}
	}

	return a.ctx.Image(), time.After(time.Millisecond * 50), nil

}

