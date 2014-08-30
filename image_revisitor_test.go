package gorevisit

import (
	"image"
	"testing"
)

func TestImageRevisitorJPEG(t *testing.T) {
	jpegTestFunc := func(src image.Image, dst image.RGBA) error {
		orig := src.Bounds()
		for x := orig.Min.X; x < orig.Max.X; x++ {
			for y := orig.Min.Y; y < orig.Max.Y; y++ {
				dst.Set(x, y, src.At(x, y))
			}
		}
		return nil
	}

	jpgMsg, err := NewRevisitMsgFromFiles("./fixtures/bob.jpg")
	if err != nil {
		t.Error(err)
	}

	if jpgMsg.ImageType() != "image/jpeg" {
		t.Error(err)
	}

	msg, err := ImageRevisitor(jpgMsg, jpegTestFunc)
	if err != nil {
		t.Error(err)
	}

	if msg.ImageType() != "image/jpeg" {
		t.Error(err)
	}
}

func TestImageRevisitorPNG(t *testing.T) {
	pngTestFunc := func(src image.Image, dst image.RGBA) error {
		orig := src.Bounds()
		for x := orig.Min.X; x < orig.Max.X; x++ {
			for y := orig.Min.Y; y < orig.Max.Y; y++ {
				dst.Set(x, y, src.At(x, y))
			}
		}
		return nil
	}

	pngMsg, err := NewRevisitMsgFromFiles("./fixtures/connie.png")
	if err != nil {
		t.Error(err)
	}

	if pngMsg.ImageType() != "image/png" {
		t.Error(err)
	}

	msg, err := ImageRevisitor(pngMsg, pngTestFunc)
	if err != nil {
		t.Error(err)
	}

	if msg.ImageType() != "image/png" {
		t.Error(err)
	}
}

func TestImageRevisitorGIF(t *testing.T) {
	gifTestFunc := func(src image.Image, dst image.RGBA) error {
		orig := src.Bounds()
		for x := orig.Min.X; x < orig.Max.X; x++ {
			for y := orig.Min.Y; y < orig.Max.Y; y++ {
				dst.Set(x, y, src.At(x, y))
			}
		}
		return nil
	}

	gifMsg, err := NewRevisitMsgFromFiles("./fixtures/bob.gif")
	if err != nil {
		t.Error(err)
	}

	if gifMsg.ImageType() != "image/gif" {
		t.Error(err)
	}

	msg, err := ImageRevisitor(gifMsg, gifTestFunc)
	if err != nil {
		t.Error(err)
	}

	if msg.ImageType() != "image/gif" {
		t.Error(err)
	}
}
