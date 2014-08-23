package main

import (
	revisit "github.com/revisitors/gorevisit"
)

func echoService(input *revisit.DecodedContent) (*revisit.DecodedContent, error) {
	return input, nil
}

func main() {
	s := revisit.NewRevisitService(echoService)
	s.Run()
}
