package gorevisit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// DecodedContent contains a type and a byte array,
// the byte array should be image data
type DecodedContent struct {
	Type string
	Data []byte
}

// Content contains a type and a string, the string
// should be a base64 encoded image
type Content struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// MetaContent contains a Content pointer
type MetaContent struct {
	Audio *Content `json:"audio"`
}

// APIMsg is a message containing Content, and MetaContent.
// the MetaContent should be audio.
type APIMsg struct {
	Content *Content     `json:"content"`
	Meta    *MetaContent `json:"meta"`
}

// NewAPIMsgFromJSON returns an APIMsg struct pointer
// from a json byte array.
func NewAPIMsgFromJSON(b []byte) (*APIMsg, error) {
	var a APIMsg
	err := json.Unmarshal(b, &a)
	return &a, err
}

// NewAPIMsgFromFiles returns an APImsg struct pointer
// given a path to an image and audio file
func NewAPIMsgFromFiles(mediaPath ...string) (*APIMsg, error) {
	if len(mediaPath) < 1 || len(mediaPath) > 2 {
		return &APIMsg{}, errors.New("must have image, may have audio")
	}

	imageBytes, err := ioutil.ReadFile(mediaPath[0])
	if err != nil {
		return &APIMsg{}, err
	}

	// FIXME: add image type detection instead of hard coded jpeg
	imageDataURI := BytesToDataURI(imageBytes, "image/jpeg")

	soundBytes, _ := ioutil.ReadFile(mediaPath[1])
	if err != nil {
		return &APIMsg{}, err
	}

	// FIXME: add sound type detection instead of hard coded ogg
	soundDataURI := BytesToDataURI(soundBytes, "audio/ogg")

	content := &Content{
		Type: "image/jpeg",
		Data: imageDataURI,
	}

	audioContent := &Content{
		Type: "audio/ogg",
		Data: soundDataURI,
	}

	metaContent := &MetaContent{
		Audio: audioContent,
	}

	apiMsg := &APIMsg{
		Content: content,
		Meta:    metaContent,
	}

	return apiMsg, nil
}

// JSON serializes a gorevisit.APIMsg back to JSON bytes
func (a *APIMsg) JSON() ([]byte, error) {
	b, err := json.Marshal(a)
	return b, err
}

// IsValid verifies that an APIMsg is valid according to the specification
func (a *APIMsg) IsValid() bool {
	// FIXME: add validation
	return true
}
