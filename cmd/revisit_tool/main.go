package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	revisit "github.com/revisitors/gorevisit"
	"net/http"
	"net/url"
	"strconv"
)

var (
	image64 string
	sound64 string
	log     = logrus.New()
)

func main() {

	var imagePath = flag.String("image", "", "path to jpeg image file")
	var soundPath = flag.String("sound", "", "path to ogg sound file")
	var output = flag.String("endpoint", "stdout", "where to output")

	flag.Parse()

	if *imagePath == "" {
		log.Fatal("--image is required")
	}

	var msg *revisit.APIMsg
	var err error

	if *soundPath != "" {
		msg, err = revisit.NewAPIMsgFromFiles(*imagePath)
	} else {
		msg, err = revisit.NewAPIMsgFromFiles(*imagePath, *soundPath)
	}

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("could not create API message")
	}

	jsonBytes, err := msg.JSON()
	if err != nil {
		log.Fatal(err)
	}

	switch *output {
	case "stdout":
		fmt.Println(string(jsonBytes))
	default:
		apiUrl, err := url.Parse(*output)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		r, _ := http.NewRequest("POST", apiUrl.String(), bytes.NewReader(jsonBytes))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("Content-Length", strconv.Itoa(len(jsonBytes)))

		resp, err := client.Do(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp.Status)
	}
}
