package main

import (
	revisit "github.com/revisitors/gorevisit"
	"image/color"
	"image/draw"
	"math/rand"
)

func noise(src draw.Image) {
	orig := src.Bounds()
	numToMod := (orig.Max.X * orig.Max.Y) / 2
	for i := 0; i < numToMod; i++ {
		x := rand.Intn(orig.Max.X)
		y := rand.Intn(orig.Max.Y)
		origColor := src.At(x, y).(color.RGBA)
		origColor.R += 30
		origColor.B += 30
		origColor.G += 30
		src.Set(x, y, origColor)
	}
}

func main() {
	// make a RevisitService instance and pass it our glitcher
	s := revisit.NewRevisitService(noise)

	// run it!
	s.Run(":8080")
}
