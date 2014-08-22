package gorevisit

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("not implemented yet")
)

type ContentInfo struct {
	Type string
	Data []byte
}

type ApiMsg struct {
	Content ContentInfo
	Meta    struct{ Audio ContentInfo }
}

func NewApiMsgFromJson(b []byte) (*ApiMsg, error) {
	return &ApiMsg{}, ErrNotImplemented
}

func (a *ApiMsg) Json() ([]byte, error) {
	return []byte(""), ErrNotImplemented
}
