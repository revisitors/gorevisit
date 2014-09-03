package gorevisit

import (
	"testing"
)

func TestNewRevisitMsgFromFilesWithJPEG(t *testing.T) {
	msg, err := NewRevisitMsgFromFiles("./fixtures/bob.jpg")
	if err != nil {
		t.Error(err)
	}

	if msg.ImageType() != "image/jpeg" {
		t.Errorf("image type should be 'image/jpeg' is %s", msg.ImageType())
	}
}

func TestNewRevisitMsgFromFilesWithGIF(t *testing.T) {
	msg, err := NewRevisitMsgFromFiles("./fixtures/bob.gif")
	if err != nil {
		t.Error(err)
	}

	if msg.ImageType() != "image/gif" {
		t.Errorf("image type should be 'image/gif' is %s", msg.ImageType())
	}
}

func TestNewRevisitMsgFromFilesWithPNG(t *testing.T) {
	msg, err := NewRevisitMsgFromFiles("./fixtures/connie.png")
	if err != nil {
		t.Error(err)
	}

	if msg.ImageType() != "image/png" {
		t.Errorf("image type should be 'image/png' is %s", msg.ImageType())
	}
}

func TestIsValidSize(t *testing.T) {
	bigImageBytes := make([]byte, 2000000)
	bigImage := &ImageData{Data: string(bigImageBytes)}
	msg := &RevisitMsg{
		Content: *bigImage,
	}
	if msg.IsValidSize() == true {
		t.Error("image should be too large but IsValidSize returned true")
	}

	smallImageBytes := make([]byte, 200)
	smallImage := &ImageData{Data: string(smallImageBytes)}
	msg = &RevisitMsg{
		Content: *smallImage,
	}
	if msg.IsValidSize() == false {
		t.Error("image should be fine but IsValidSize returned false")
	}

}
