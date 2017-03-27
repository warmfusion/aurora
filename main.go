package main

import (
	"flag"
	"image"
	"image/color"
	"math/rand"
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
	maxSparks  = flag.Int("max-sparks", 100, "sparks (0-100)")
	sparks     = make([]Spark, *maxSparks) // the slice v now refers to a new array of 100 ints
)

type Spark struct {
	radius    float64
	color     color.Color
	position  SparkPoint
	direction SparkPoint
}

type SparkPoint struct {
	X float64
	Y float64
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
	ctx          *gg.Context
	position     image.Point
	dir          image.Point
	frameCounter int
	stroke       int
}

func NewAnimation(sz image.Point) *Animation {
	for i := range sparks {
		sparks[i] = Spark{
			position:  SparkPoint{32, 16},
			radius:    3.0,
			color:     colorful.Hsv(rand.Float64()*360, 0.56, 0.3),
			direction: SparkPoint{(rand.Float64() * 2) - 1, (rand.Float64() * 2) - 1},
		}
	}
	return &Animation{
		ctx:          gg.NewContext(sz.X, sz.Y),
		dir:          image.Point{1, 1},
		frameCounter: 0,
		stroke:       6,
	}
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {

	a.frameCounter++
	c := a.frameCounter % 360

	// Restart
	a.ctx.SetColor(color.Black)
	a.ctx.Clear()

	for s := range sparks {

		spark := &sparks[s]

		// Linear movement for now
		spark.position.X += spark.direction.X
		spark.position.Y += spark.direction.Y

		a.ctx.SetColor(spark.color)
		a.ctx.DrawCircle(spark.position.X, spark.position.Y, float64(spark.radius))
		a.ctx.Fill()
	}
	/*
		a.ctx.SetColor(colorful.Hsv(float64(c), 0.56, 0.3))
		a.ctx.DrawCircle(25, 5.0, float64(a.stroke))
		a.ctx.Fill()

		a.ctx.SetColor(colorful.Hsv(float64(c+45), 0.56, 0.3))
		a.ctx.DrawCircle(30, 10.0, float64(a.stroke))
		a.ctx.Fill()

		a.ctx.SetColor(colorful.Hsv(float64((c+90)%360), 0.56, 0.3))
		a.ctx.DrawCircle(35, 15.0, float64(a.stroke))
		a.ctx.Fill()

		a.ctx.SetColor(colorful.Hsv(float64((c+135)%360), 0.56, 0.3))
		a.ctx.DrawCircle(40, 20.0, float64(a.stroke))
		a.ctx.Fill()

		a.ctx.SetColor(colorful.Hsv(float64((c+180)%360), 0.56, 0.3))
		a.ctx.DrawCircle(45, 25.0, float64(a.stroke))
		a.ctx.Fill()
	*/
	return a.ctx.Image(), time.After(time.Millisecond * 50), nil
}
