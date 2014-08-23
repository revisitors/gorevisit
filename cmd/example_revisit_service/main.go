package main

import (
	revisit "github.com/revisitors/gorevisit"
)

func echoService(input *revisit.APIMsg) (*revisit.APIMsg, error) {
	return input, nil
}

func main() {
	s := revisit.NewRevisitService(echoService)
	s.Run()
}
