package gorevisit

import (
	"io/ioutil"
	"net/http"
)

const (
	PayloadLimit = 1000000
)

type RevisitService struct {
	// INTERFACE
}

func (rs *RevisitService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		rs.HandlePost(w, r)
	default:
		http.Error(w, "ROFLMAO", http.StatusMethodNotAllowed)
		return
	}
}

func (rs *RevisitService) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "ROFLMAO", http.StatusUnsupportedMediaType)
		return
	}

	payload := http.MaxBytesReader(w, r.Body, PayloadLimit)
	payloadBytes, err := ioutil.ReadAll(payload)
	if err != nil {
		http.Error(w, "ROFLMAO", http.StatusRequestEntityTooLarge)
		return
	}

	apiMsg, err := NewAPIMsgFromJSON(payloadBytes)
	if err != nil {
		http.Error(w, "ROFLMAO", http.StatusUnsupportedMediaType)
		return
	}

	if !apiMsg.IsValid() {
		http.Error(w, "ROFLMAO", http.StatusUnsupportedMediaType)
		return
	}
}
