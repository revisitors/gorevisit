package gorevisit

import (
	"encoding/base64"
	"io/ioutil"
	"testing"
)

func getValidTestAPIMsg(t *testing.T) *APIMsg {
	imageBytes, _ := ioutil.ReadFile("./fixtures/bob.jpg")
	image64 := base64.StdEncoding.EncodeToString(imageBytes)

	soundBytes, _ := ioutil.ReadFile("./fixtures/scream.ogg")
	sound64 := base64.StdEncoding.EncodeToString(soundBytes)

	content := &Content{
		Type: "image/jpeg",
		Data: image64,
	}

	audioContent := &Content{
		Type: "audio/ogg",
		Data: sound64,
	}

	metaContent := &MetaContent{
		Audio: audioContent,
	}

	apiMsg := &APIMsg{
		Content: content,
		Meta:    metaContent,
	}

	return apiMsg

}

func TestJSON(t *testing.T) {
	msg := getValidTestAPIMsg(t)
	_, err := msg.JSON()
	if err != nil {
		t.Error(err)
	}
}

func TestNewAPIMsgFromJSON(t *testing.T) {
	msg := getValidTestAPIMsg(t)
	jsonMsg, _ := msg.JSON()

	newMsg, err := NewAPIMsgFromJSON(jsonMsg)
	if err != nil {
		t.Error(err)
	}

	if !newMsg.IsValid() {
		t.Error("message is not valid")
	}
}
