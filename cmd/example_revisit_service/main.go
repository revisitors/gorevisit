package main

import (
	"github.com/gorilla/mux"
	revisit "github.com/revisitors/gorevisit"
	"net/http"
)

func echoService(input *revisit.APIMsg) (*revisit.APIMsg, error) {
	return input, nil
}

func simpleBlend(input *revisit.APIMsg) (*revisit.APIMsg, error) {
	imageContent, err := revisit.DataURIToDecodedContent(input.Content.Data)
	if err != nil {
		return input, err
	}

	// TODO: add transformation
	newImageBytes := imageContent.Data

	// FIXME: fix hard coded image type
	input.Content.Data = revisit.BytesToDataURI(newImageBytes, "image/jpeg")

	return input, nil
}

func main() {
	s := revisit.NewRevisitService(echoService)
	r := mux.NewRouter()
	r.HandleFunc("/", s.ServeHTTP)
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
