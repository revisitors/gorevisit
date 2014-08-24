package gorevisit

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"image"
	"io/ioutil"
	"net/http"
)

const (
	payloadLimit = 1e6
)

var (
	log = logrus.New()
)

type RevisitService struct {
	imageTransformer func(image.Image, image.RGBA) error
}

func NewRevisitService(it func(image.Image, image.RGBA) error) *RevisitService {
	return &RevisitService{imageTransformer: it}
}

func (rs *RevisitService) ServiceCheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "HEAD":
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (rs *RevisitService) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("%v", r)
	switch r.Method {
	case "POST":
		rs.PostHandler(w, r)
	case "HEAD":
		w.WriteHeader(http.StatusOK)
		return
	default:
		log.Infof("%v", r.Method)
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func (rs *RevisitService) PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	payloadReadCloser := http.MaxBytesReader(w, r.Body, payloadLimit)
	payloadBytes, err := ioutil.ReadAll(payloadReadCloser)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusRequestEntityTooLarge)
		return
	}

	var msg *RevisitMsg
	decoder := json.NewDecoder(bytes.NewReader(payloadBytes))
	err = decoder.Decode(&msg)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	newMsg, err := ImageRevisitor(msg, rs.imageTransformer)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(newMsg)
	return
}

func (rs *RevisitService) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", rs.ServiceCheckHandler)
	r.HandleFunc("/service", rs.ServiceHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
