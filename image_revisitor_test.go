package gorevisit

import (
	"image"
	"testing"
)

func testFunc(src image.Image, dst image.RGBA) error {
	orig := src.Bounds()
	for x := orig.Min.X; x < orig.Max.X; x++ {
		for y := orig.Min.Y; y < orig.Max.Y; y++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
	return nil
}

func TestImageRevisitorJPEG(t *testing.T) {
}
