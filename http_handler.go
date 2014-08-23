package gorevisit

import (
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const (
	payloadLimit = 1e6
)

var (
	log = logrus.New()
)

// RevisitService holds context for a POST handler for revisit
type RevisitService struct {
	Transform func(*APIMsg) (*APIMsg, error)
}

// NewRevisitService constructs a new Revisit service given a transform function
func NewRevisitService(t func(*APIMsg) (*APIMsg, error)) *RevisitService {
	return &RevisitService{Transform: t}
}

// ServeHTTP implements a Revisit service to be passed to a mux
func (rs *RevisitService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// PostHandler handles a POST to a revisit service
func (rs *RevisitService) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	payload := http.MaxBytesReader(w, r.Body, payloadLimit)
	payloadBytes, err := ioutil.ReadAll(payload)
	if err != nil {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusRequestEntityTooLarge)
		return
	}

	log.WithFields(logrus.Fields{"type": "request"}).Info(string(payloadBytes))

	original, err := NewAPIMsgFromJSON(payloadBytes)
	if err != nil {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	if !original.IsValid() {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	transformed, err := rs.Transform(original)

	if err != nil {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
			"error":  err,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	if !transformed.IsValid() {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	transformedJSON, err := transformed.JSON()

	if err != nil {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
			"error":  err,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(transformedJSON)

	return
}
