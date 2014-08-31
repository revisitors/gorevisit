package gorevisit

import (
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"testing"
)

func mockTransform(src draw.Image) {
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

func TestNewRevisitImageWithJPEG(t *testing.T) {
	jpegMsg, err := NewRevisitMsgFromFiles("./fixtures/bob.jpg")
	if err != nil {
		t.Fatal(err)
	}

	ri, err := NewRevisitImageFromMsg(jpegMsg)
	if err != nil {
		t.Fatal(err)
	}

	if len(ri.rgbas) != 1 {
		t.Errorf("ri.rgbas length should be 1, is %d", len(ri.rgbas))
	}

	if len(ri.delay) != 1 {
		t.Errorf("ri.delay length should be 1, is %d", len(ri.delay))
	}

	if ri.loopCount != 0 {
		t.Errorf("loopCount should be 0, is %d", ri.loopCount)
	}

	ri.Transform(mockTransform)

	m, err := ri.RevisitMsg()
	if err != nil {
		t.Error(err)
	}

	_, err = jpeg.Decode(m.ImageByteReader())
	if err != nil {
		t.Error(err)
	}
}

func TestNewRevisitImageWithPNG(t *testing.T) {
	pngMsg, err := NewRevisitMsgFromFiles("./fixtures/connie.png")
	if err != nil {
		t.Fatal(err)
	}

	ri, err := NewRevisitImageFromMsg(pngMsg)
	if err != nil {
		t.Fatal(err)
	}

	if len(ri.rgbas) != 1 {
		t.Errorf("ri.rgbas length should be 1, is %d", len(ri.rgbas))
	}

	if len(ri.delay) != 1 {
		t.Errorf("ri.delay length should be 1, is %d", len(ri.delay))
	}

	if ri.loopCount != 0 {
		t.Errorf("loopCount should be 0, is %d", ri.loopCount)
	}

	ri.Transform(mockTransform)

	m, err := ri.RevisitMsg()
	if err != nil {
		t.Error(err)
	}

	_, err = png.Decode(m.ImageByteReader())
	if err != nil {
		t.Error(err)
	}
}

func TestNewRevisitImageWithGIF(t *testing.T) {
	gifMsg, err := NewRevisitMsgFromFiles("./fixtures/bob.gif")
	if err != nil {
		t.Fatal(err)
	}

	ri, err := NewRevisitImageFromMsg(gifMsg)
	if err != nil {
		t.Fatal(err)
	}

	if len(ri.rgbas) != 4 {
		t.Errorf("ri.rgbas length should be 4, is %d", len(ri.rgbas))
	}

	if len(ri.delay) != 4 {
		t.Errorf("ri.delay length should be 4, is %d", len(ri.delay))
	}

	if ri.loopCount != 0 {
		t.Errorf("loopCount should be 0, is %d", ri.loopCount)
	}

	ri.Transform(mockTransform)

	m, err := ri.RevisitMsg()
	if err != nil {
		t.Error(err)
	}

	_, err = gif.Decode(m.ImageByteReader())
	if err != nil {
		t.Error(err)
	}
}
