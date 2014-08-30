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

	if !ri.IsJPEG() {
		t.Error("isJPEG returned false should be true")
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

	if !ri.IsPNG() {
		t.Error("isPNG returned false should be true")
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

	if !ri.IsGIF() {
		t.Error("isGIF returned false should be true")
	}
}
