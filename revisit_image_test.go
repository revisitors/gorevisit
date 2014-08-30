package gorevisit

import (
	"testing"
)

func TestNewRevisitImageWithJPEG(t *testing.T) {
	jpegMsg, err := NewRevisitMsgFromFiles("./fixtures/bob.jpg")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewRevisitImageFromMsg(jpegMsg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRevisitImageWithPNG(t *testing.T) {
	pngMsg, err := NewRevisitMsgFromFiles("./fixtures/connie.png")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewRevisitImageFromMsg(pngMsg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRevisitImageWithGIF(t *testing.T) {
	gifMsg, err := NewRevisitMsgFromFiles("./fixtures/bob.gif")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewRevisitImageFromMsg(gifMsg)
	if err != nil {
		t.Fatal(err)
	}
}
