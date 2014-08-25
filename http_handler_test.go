package gorevisit

import (
	"bytes"
	"encoding/json"
	"image"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func echoService(src image.Image, dst image.RGBA) error {
	orig := src.Bounds()
	for x := orig.Min.X; x < orig.Max.X; x++ {
		for y := orig.Min.Y; y < orig.Max.Y; y++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
	return nil
}

func TestRevisitHandlerPost(t *testing.T) {
	msg, err := NewRevisitMsgFromFiles("./fixtures/bob.jpg", "./fixtures/scream.ogg")
	if err != nil {
		t.Fatal(err)
	}

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("POST", "http://whatever", bytes.NewReader(jsonBytes))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonBytes)))

	w := httptest.NewRecorder()

	service := NewRevisitService(echoService)
	service.serviceHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d status, received %d", http.StatusOK, w.Code)
	}
}
