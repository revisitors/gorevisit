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
	// revisit spec size limit is 1 MB per media item.
	// this limits payload size of entire message to
	// 2 meg with a little overhead for json
	payloadLimit = 2001000
)

var (
	log = logrus.New()
)

// RevisitService holds the necessary context for a Revisit.link service.
// Currently gorevisit only handles image data.
type RevisitService struct {
	glitcher func(draw.Image)
}

// NewRevisitService given an image transformation function, returns
// a new Revisit.link service.
// The image transformation service receives a draw.Image interface
// as an argument.  Note that draw.Image also implements image.Image.
// For details see:
// * http://golang.org/pkg/image/draw/
// * http://golang.org/pkg/image/#Image
func NewRevisitService(g func(draw.Image)) *RevisitService {
	return &RevisitService{glitcher: g}
}

// headHandler responds to availability requests from a Revisit.link hub
// see: http://revisit.link/spec.html
func (rs *RevisitService) headHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("%v", r)
	w.WriteHeader(http.StatusOK)
}

// postHandler accepts POSTed revisit messages from a Revisit.link hub,
// transforms the message, and returns the transformed message to the hub.
// See: http://revisit.link/spec.html
func (rs *RevisitService) postHandler(w http.ResponseWriter, r *http.Request) {

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

	// construct a RevisitImage from the payload
	ri, err := NewRevisitImageFromMsg(msg)
	if err != nil {
		log.Errorf("error decoding json: %d", http.StatusUnsupportedMediaType)
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	// for each image frame (only 1 if jpeg or png), call the glitcher
	for _, rgba := range ri.rgbas {
		rs.glitcher(draw.Image(&rgba))
	}

	// create a new message from the modified image
	newMsg, err := ri.RevisitMsg()
	if err != nil {
		log.Errorf("error decoding json: %d", http.StatusInternalServerError)
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(newMsg)
}

// Run starts the Revisit.link service
func (rs *RevisitService) Run(port string) {
	rmux := mux.NewRouter()

	rmux.HandleFunc("/", rs.headHandler)

	rmux.HandleFunc("/service", rs.postHandler).
		Methods("POST").
		Headers("Content-Type", "application/json")

	http.Handle("/", rmux)
	log.Infof("Listening to %s", port)
	http.ListenAndServe(port, rmux)
}
