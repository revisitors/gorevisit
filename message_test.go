package gorevisit

import (
	"testing"
)

var (
	// FIXME: replace with actual test fixture
	testMsg = []byte("")
)

func TestNewApiMsgFromJson(t *testing.T) {
	_, err := NewApiMsgFromJson(testMsg)
	if err != nil {
		t.Error(err)
	}
}

func TestJson(t *testing.T) {
	msg, err := NewApiMsgFromJson(testMsg)
	if err != nil {
		t.Error(err)
	}

	_, err = msg.Json()
	if err != nil {
		t.Error(err)
	}
}
