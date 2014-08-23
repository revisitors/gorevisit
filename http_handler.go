package gorevisit

import (
	"io/ioutil"
	"net/http"
)

const (
	PayloadLimit = 1000000
)

type RevisitService struct {
	transform func(*APIMsg) (*APIMsg, error)
}

func NewRevisitService(t func(*APIMsg) (*APIMsg, error)) *RevisitService {
	return &RevisitService{transform: t}
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

	original, err := NewAPIMsgFromJSON(payloadBytes)
	if err != nil {
		http.Error(w, "ROFLMAO", http.StatusUnsupportedMediaType)
		return
	}

	if !original.IsValid() {
		http.Error(w, "ROFLMAO", http.StatusUnsupportedMediaType)
		return
	}

	morphed, err := rs.transform(original)

	if err != nil {
		http.Error(w, "ROFLMAO", http.StatusInternalServerError)
		return
	}

	if !morphed.IsValid() {
		http.Error(w, "ROFLMAO", http.StatusInternalServerError)
		return
	}

	morphedJson, err := morphed.JSON()

	if err != nil {
		http.Error(w, "ROFLMAO", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(morphedJson)

	return
}
