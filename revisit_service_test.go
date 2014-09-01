package gorevisit

import (
	"bytes"
	"encoding/json"
	"image/color"
	"image/draw"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func noise(src draw.Image) {
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

	service := NewRevisitService(noise)
	service.serviceHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d status, received %d", http.StatusOK, w.Code)
	}
}
