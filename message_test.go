package gorevisit

import (
	"testing"
)

func getValidTestAPIMsg(t *testing.T) *APIMsg {
	msg, err := NewAPIMsgFromFiles("./fixtures/bob.jpg", "./fixtures/scream.ogg")
	if err != nil {
		t.Fatal(err)
	}
	return msg
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
