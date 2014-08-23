package gorevisit

import (
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
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
	Transform func(*DecodedContent) (*DecodedContent, error)
}

// NewRevisitService constructs a new Revisit service given a transform function
func NewRevisitService(t func(*DecodedContent) (*DecodedContent, error)) *RevisitService {
	return &RevisitService{Transform: t}
}

// ServiceCheckHandler handles presence checks from the hub
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

// TransformationHandler implements a Revisit service to be passed to a mux
func (rs *RevisitService) TransformationHandler(w http.ResponseWriter, r *http.Request) {
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
			"status": http.StatusRequestEntityTooLarge,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusRequestEntityTooLarge)
		return
	}

	log.WithFields(logrus.Fields{"type": "request"}).Info(string(payloadBytes))

	apiMsg, err := NewAPIMsgFromJSON(payloadBytes)
	if err != nil {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	if !apiMsg.IsValid() {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	inDecodedContent, err := apiMsg.GetImageDecodedContent()
	if err != nil {
		log.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError,
			"error":  err,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	outDecodedContent, err := rs.Transform(inDecodedContent)
	if err != nil {
		log.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError,
			"error":  err,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	newBase64 := BytesToDataURI(outDecodedContent.Data, outDecodedContent.Type)
	apiMsg.Content.Data = newBase64

	if !apiMsg.IsValid() {
		log.WithFields(logrus.Fields{
			"status": http.StatusUnsupportedMediaType,
		}).Error("HTTP Error")

		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	transformedJSON, err := apiMsg.JSON()

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

// Run starts the service
func (rs *RevisitService) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", rs.ServiceCheckHandler)
	r.HandleFunc("/service", rs.TransformationHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
