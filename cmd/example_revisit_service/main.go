package main

import (
	revisit "github.com/revisitors/gorevisit"
	"image"
	"image/color"
	"log"
	"math/rand"
)

func noise(src image.Image, dst image.RGBA) error {
	orig := src.Bounds()

	for x := orig.Min.X; x < orig.Max.X; x++ {
		for y := orig.Min.Y; y < orig.Max.Y; y++ {
			dst.Set(x, y, src.At(x, y))
		}
	}

	numToMod := (orig.Max.X * orig.Max.Y) / 2
	for i := 0; i < numToMod; i++ {
		x := rand.Intn(orig.Max.X)
		y := rand.Intn(orig.Max.Y)
		if prev, ok := src.At(x, y).(color.RGBA); ok {
			prev.R += 30
			prev.B -= 30
			prev.G += 30
			dst.Set(x, y, prev)
		}
		if prev, ok := src.At(x, y).(color.YCbCr); ok {
			prev.Cr += 30
			prev.Cb -= 30
			prev.Y += 30
			dst.Set(x, y, prev)
		}
	}

	return nil

}

func main() {
	log.Println("starting")
	s := revisit.NewRevisitService(noise)
	s.Run()
}
