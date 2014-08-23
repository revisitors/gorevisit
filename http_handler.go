package gorevisit

import (
	"io/ioutil"
	"net/http"
)

const (
	payloadLimit = 1000000
)

// RevisitService holds context for a POST handler for revisit
type RevisitService struct {
	transform func(*APIMsg) (*APIMsg, error)
}

// NewRevisitService constructs a new Revisit service given a transform function
func NewRevisitService(t func(*APIMsg) (*APIMsg, error)) *RevisitService {
	return &RevisitService{transform: t}
}

// ServeHTTP implements a Revisit service to be passed to a mux
func (rs *RevisitService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		rs.PostHandler(w, r)
	default:
		http.Error(w, "ROTFL", http.StatusMethodNotAllowed)
		return
	}
}

// PostHandler handles a POST to a revisit service
func (rs *RevisitService) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	payload := http.MaxBytesReader(w, r.Body, payloadLimit)
	payloadBytes, err := ioutil.ReadAll(payload)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusRequestEntityTooLarge)
		return
	}

	original, err := NewAPIMsgFromJSON(payloadBytes)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	if !original.IsValid() {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	morphed, err := rs.transform(original)

	if err != nil {
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	if !morphed.IsValid() {
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	morphedJson, err := morphed.JSON()

	if err != nil {
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(morphedJson)

	return
}
