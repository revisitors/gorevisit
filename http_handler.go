package gorevisit

import (
	"bytes"
	"encoding/base64"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	payloadLimit = 1e6
)

var (
	log = logrus.New()
)

// RevisitService holds context for a POST handler for revisit
type RevisitService struct {
	Transform func(image.Image) (image.Image, error)
}

// NewRevisitService constructs a new Revisit service given a transform function
func NewRevisitService(t func(image.Image) (image.Image, error)) *RevisitService {
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

	// Check for appropriate header
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	// Check for appropriate message size
	payload := http.MaxBytesReader(w, r.Body, payloadLimit)
	payloadBytes, err := ioutil.ReadAll(payload)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusRequestEntityTooLarge)
		return
	}

	log.WithFields(logrus.Fields{"type": "request"}).Info(string(payloadBytes))

	// Convery POSTed JSON to apiMSG
	apiMsg, err := NewAPIMsgFromJSON(payloadBytes)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	// Check validity of apiMsg
	if !apiMsg.IsValid() {
		http.Error(w, "ROTFL", http.StatusUnsupportedMediaType)
		return
	}

	// get []byte array of image content and string of type
	dataURI := apiMsg.Content.Data
	parts := strings.Split(dataURI, ",")
	contentBytes, err := base64.StdEncoding.DecodeString(parts[1])

	// convert byte array of img data to an actual image and get it's type
	inImg, inImgType, err := image.Decode(bytes.NewBuffer(contentBytes))
	if err != nil {
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}
	inImgType = "image/" + inImgType

	// transform the actual image
	_, err = rs.Transform(inImg)
	if err != nil {
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	// create a buffer to write the encoded image to
	imgBuf := new(bytes.Buffer)

	// encode the image back to bytes
	switch inImgType {
	case "image/gif":
		log.Info("encoding gif")
		gif.Encode(imgBuf, inImg, nil)
	case "image/jpeg":
		log.Info("encoding jpeg")
		jpeg.Encode(imgBuf, inImg, nil)
	case "image/png":
		log.Info("encoding png")
		png.Encode(imgBuf, inImg)
	}

	newBase64 := BytesToDataURI(imgBuf.Bytes(), inImgType)
	apiMsg.Content.Data = newBase64

	if !apiMsg.IsValid() {
		http.Error(w, "ROTFL", http.StatusInternalServerError)
		return
	}

	transformedJSON, err := apiMsg.JSON()

	if err != nil {
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
