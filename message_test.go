package gorevisit

import (
	"encoding/base64"
	"io/ioutil"
	"testing"
)

func getValidTestMessage(t *testing.T) *ApiMsg {
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

	apiMsg := &ApiMsg{
		Content: content,
		Meta:    metaContent,
	}

	return apiMsg

}

func TestJson(t *testing.T) {
	msg := getValidTestMessage(t)
	_, err := msg.Json()
	if err != nil {
		t.Error(err)
	}
}

func TestNewApiMsgFromJson(t *testing.T) {
	msg := getValidTestMessage(t)
	jsonMsg, _ := msg.Json()
	_, err := NewApiMsgFromJson(jsonMsg)
	if err != nil {
		t.Error(err)
	}
}
