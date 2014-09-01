package gorevisit

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"image/draw"
	"io/ioutil"
	"net/http"
)

const (
	payloadLimit = 1e6
)

var (
	log = logrus.New()
)

// RevisitService holds the necessary context for a Revisit.link service.
// Currently, this consists of an imageTransformer
type RevisitService struct {
	imageTransformer func(draw.Image)
}

// NewRevisitService given an image transformation function, returns
// a new Revisit.link service
func NewRevisitService(it func(draw.Image)) *RevisitService {
	return &RevisitService{imageTransformer: it}
}

// serviceCheckHandler responds to availability requests from a Revisit.link hub
func (rs *RevisitService) serviceCheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "HEAD":
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// serviceHandler appropriately routes service requests from a Revisit.link hub
func (rs *RevisitService) serviceHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("%v", r)

	switch r.Method {
	case "POST":
		rs.postHandler(w, r)
	case "HEAD":
		w.WriteHeader(http.StatusOK)
		return
	default:
		log.Infof("%v", r.Method)
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

// postHandler accepts POSTed revisit messages from a Revisit.link hub,
// transforms the message, and returns the transformed message to the hub
func (rs *RevisitService) postHandler(w http.ResponseWriter, r *http.Request) {

	// check for valid header
	if r.Header.Get("Content-Type") != "application/json" {
		log.Errorf("error invalid header: %d", http.StatusUnsupportedMediaType)
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	// make sure message isn't too large
	payloadReadCloser := http.MaxBytesReader(w, r.Body, payloadLimit)
	payloadBytes, err := ioutil.ReadAll(payloadReadCloser)
	if err != nil {
		log.Errorf("error reading payload: %d", http.StatusRequestEntityTooLarge)
		http.Error(w, "ROTFL", http.StatusRequestEntityTooLarge)
		return
	}

	// decode the payload into a RevisitMsg
	var msg *RevisitMsg
	decoder := json.NewDecoder(bytes.NewReader(payloadBytes))
	err = decoder.Decode(&msg)
	if err != nil {
		log.Errorf("error decoding json: %d", http.StatusUnsupportedMediaType)
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	ri, err := NewRevisitImageFromMsg(msg)
	if err != nil {
		log.Errorf("error decoding json: %d", http.StatusUnsupportedMediaType)
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	for _, rgba := range ri.rgbas {
		rs.imageTransformer(draw.Image(&rgba))
	}

	newMsg, err := ri.RevisitMsg()
	if err != nil {
		log.Errorf("error decoding json: %d", http.StatusInternalServerError)
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(newMsg)
	return
}

// Run starts the Revisit.link service
func (rs *RevisitService) Run(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/", rs.serviceCheckHandler)
	r.HandleFunc("/service", rs.serviceHandler)
	http.Handle("/", r)
	log.Infof("Listening to %s", port)
	http.ListenAndServe(port, r)
}
