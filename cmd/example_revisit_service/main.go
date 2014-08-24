package main

import (
	revisit "github.com/revisitors/gorevisit"
	"image"
	"image/color"
	"log"
)

type ImageSet interface {
	Set(x, y int, c color.Color)
}

// rotate stolen from Justin Abrahms while I make sure my interface works
// and is compatible!  See - https://github.com/justinabrahms/rotate.revisit.link
func rotate(m image.Image, dst image.RGBA) error {
	orig := m.Bounds()
	for x := orig.Min.X; x < orig.Max.X; x++ {
		for y := orig.Min.Y; y < orig.Max.Y; y++ {
			dst.Set(y,
				orig.Max.Y-x+1,
				m.At(x, y))
		}
	}
	return nil
}

func main() {
	log.Println("starting")
	s := revisit.NewRevisitService(rotate)
	s.Run()
}
