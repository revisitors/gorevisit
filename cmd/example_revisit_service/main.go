package main

import (
	revisit "github.com/revisitors/gorevisit"
	"image"
)

func echoService(img image.Image) (image.Image, error) {
	return img, nil
}

func main() {
	s := revisit.NewRevisitService(echoService)
	s.Run()
}
