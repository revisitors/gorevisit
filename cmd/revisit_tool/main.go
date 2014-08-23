package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	revisit "github.com/revisitors/go.revisit.link"
	"io/ioutil"
	"log"
)

var (
	image64 string
	sound64 string
)

func main() {
	var imageUrl = flag.String("image", "", "path to jpeg image file")
	var soundUrl = flag.String("sound", "", "path to ogg sound file")
	flag.Parse()

	if *imageUrl != "" {
		imageHeader := "data:image/jpeg;base64,"
		imageBytes, err := ioutil.ReadFile(*imageUrl)
		if err != nil {
			log.Fatal(err)
		}
		image64 = imageHeader + base64.StdEncoding.EncodeToString(imageBytes)
	}

	if *soundUrl != "" {
		soundHeader := "data:audio/ogg;base64,"
		soundBytes, err := ioutil.ReadFile(*soundUrl)
		if err != nil {
			log.Fatal(err)
		}
		sound64 = soundHeader + base64.StdEncoding.EncodeToString(soundBytes)
	}

	content := &revisit.Content{
		Type: "image/jpeg",
		Data: image64,
	}

	audioContent := &revisit.Content{
		Type: "audio/ogg",
		Data: sound64,
	}

	metaContent := &revisit.MetaContent{
		Audio: audioContent,
	}

	apiMsg := &revisit.ApiMsg{
		Content: content,
		Meta:    metaContent,
	}

	jsonBytes, err := apiMsg.Json()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))

}
