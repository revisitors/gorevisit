package gorevisit

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotImplemented = errors.New("not implemented yet")
)

// Content contains a type and a byte array
// and should be an image
type Content struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// AudioContent contains a Content pointer
type MetaContent struct {
	Audio *Content `json:"audio"`
}

// ApiMsg is a message containing Content, and MetaContent.
// the MetaContent should be audio.
type ApiMsg struct {
	Content *Content     `json:"content"`
	Meta    *MetaContent `json:"meta"`
}

// NewApiMsgFromJson returns an ApiMsg struct pointer
// from a json byte array.
func NewApiMsgFromJson(b []byte) (*ApiMsg, error) {
	var a ApiMsg
	err := json.Unmarshal(b, &a)
	return &a, err
}

// Json serializes a gorevisit.ApiMsg back to JSON bytes
func (a *ApiMsg) Json() ([]byte, error) {
	b, err := json.Marshal(a)
	return b, err
}
