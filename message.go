package gorevisit

import (
	"encoding/json"
)

// Content contains a type and a byte array
// and should be an image
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

// JSON serializes a gorevisit.APIMsg back to JSON bytes
func (a *APIMsg) JSON() ([]byte, error) {
	b, err := json.Marshal(a)
	return b, err
}
