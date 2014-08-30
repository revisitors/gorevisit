package gorevisit

import (
	"testing"
)

func TestNewRevisitImageWithJPEG(t *testing.T) {
	jpegMsg, err := NewRevisitMsgFromFiles("./fixtures/bob.jpg")
	if err != nil {
		t.Fatal(err)
	}

	ri, err := NewRevisitImageFromMsg(jpegMsg)
	if err != nil {
		t.Fatal(err)
	}

	if len(ri.images) != 1 {
		t.Errorf("ri.images length should be 1, is %d", len(ri.images))
	}

	if len(ri.palettes) != 0 {
		t.Errorf("ri.palettes length should be 0, is %d", len(ri.palettes))
	}

	if len(ri.delay) != 1 {
		t.Errorf("ri.delay length should be 1, is %d", len(ri.delay))
	}

	if ri.loopCount != 0 {
		t.Errorf("loopCount should be 0, is %d", ri.loopCount)
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

	if len(ri.images) != 1 {
		t.Errorf("ri.images length should be 1, is %d", len(ri.images))
	}

	if len(ri.palettes) != 0 {
		t.Errorf("ri.palettes length should be 0, is %d", len(ri.palettes))
	}

	if len(ri.delay) != 1 {
		t.Errorf("ri.delay length should be 1, is %d", len(ri.delay))
	}

	if ri.loopCount != 0 {
		t.Errorf("loopCount should be 0, is %d", ri.loopCount)
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

	if len(ri.images) != 0 {
		t.Errorf("ri.images length should be 0, is %d", len(ri.images))
	}

	if len(ri.palettes) != 4 {
		t.Errorf("ri.palettes length should be 4, is %d", len(ri.palettes))
	}

	if len(ri.delay) != 4 {
		t.Errorf("ri.delay length should be 4, is %d", len(ri.delay))
	}

	if ri.loopCount != 0 {
		t.Errorf("loopCount should be 0, is %d", ri.loopCount)
	}

}
