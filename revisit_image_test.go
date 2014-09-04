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

	if len(ri.Rgbas) != 1 {
		t.Errorf("ri.Rgbas length should be 1, is %d", len(ri.Rgbas))
	}

	if len(ri.Delay) != 1 {
		t.Errorf("ri.Delay length should be 1, is %d", len(ri.Delay))
	}

	if ri.LoopCount != 0 {
		t.Errorf("LoopCount should be 0, is %d", ri.LoopCount)
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

	if len(ri.Rgbas) != 1 {
		t.Errorf("ri.Rgbas length should be 1, is %d", len(ri.Rgbas))
	}

	if len(ri.Delay) != 1 {
		t.Errorf("ri.Delay length should be 1, is %d", len(ri.Delay))
	}

	if ri.LoopCount != 0 {
		t.Errorf("LoopCount should be 0, is %d", ri.LoopCount)
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

	if len(ri.Rgbas) != 4 {
		t.Errorf("ri.Rgbas length should be 4, is %d", len(ri.Rgbas))
	}

	if len(ri.Delay) != 4 {
		t.Errorf("ri.Delay length should be 4, is %d", len(ri.Delay))
	}

	if ri.LoopCount != 0 {
		t.Errorf("LoopCount should be 0, is %d", ri.LoopCount)
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
