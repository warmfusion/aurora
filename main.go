package main

import (
	"flag"
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

var (
	rows       = flag.Int("led-rows", 32, "number of rows supported")
	parallel   = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain      = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness = flag.Int("brightness", 100, "brightness (0-100)")
)

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
	steps    int
	stroke   int
}

func NewAnimation(sz image.Point) *Animation {
	return &Animation{
		ctx:    gg.NewContext(sz.X, sz.Y),
		dir:    image.Point{1, 1},
		steps:  0,
		stroke: 10,
	}
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {

	a.ctx.SetColor(color.Black)
	a.ctx.Clear()

	a.steps += 20
	c := a.steps % 360

	a.ctx.DrawCircle(30, 10.0, float64(a.stroke))
	a.ctx.SetColor(colorful.Hsv(float64(c), 0.56, 0.3))
	a.ctx.Fill()

	a.ctx.DrawCircle(35, 15.0, float64(a.stroke))
	a.ctx.SetColor(colorful.Hsv(float64((c+90)%360), 0.56, 0.3))
	a.ctx.Fill()

	a.ctx.DrawCircle(40, 20.0, float64(a.stroke))
	a.ctx.SetColor(colorful.Hsv(float64((c+180)%360), 0.56, 0.3))
	a.ctx.Fill()

	return a.ctx.Image(), time.After(time.Millisecond * 50), nil
}
