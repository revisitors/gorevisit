/*
Portions of this code from: https://github.com/Knorkebrot/aimg/
ISC license
Copyright (c) 2014, Bjoern Oelke <bo@kbct.de>
Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.
THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
*/

package main

import (
	"fmt"
	"github.com/Knorkebrot/ansirgb"
	"github.com/gorilla/websocket"
	"github.com/revisitors/gorevisit"
	"image"
	"log"
	"net"
	"net/http"
	"net/url"
)

var (
	width = 80
)

type block struct {
	top    *ansirgb.Color
	bottom *ansirgb.Color
}

func cursorUp(count int) {
	fmt.Printf("\033[%dA", count)
}
func reset() {
	// add a space to prevent artifacts after resizing
	fmt.Printf("\033[0m ")
}

func main() {
	conn, err := net.Dial("tcp", "ws.revisit.link:80")
	if err != nil {
		panic(err)
	}

	h := make(http.Header)
	u, _ := url.Parse("http://ws.revisit.link:80/message")
	ws, _, err := websocket.NewClient(conn, u, h, 1024, 1024)
	if err != nil {
		log.Println(err)
	}

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		id := gorevisit.ImageData{
			Data: string(p),
		}

		msg := &gorevisit.RevisitMsg{
			Content: id,
		}

		reader := msg.Content.ByteReader()
		img, _, err := image.Decode(reader)
		if err != nil {
			log.Println(err)
		}

		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()
		if imgWidth < width {
			width = imgWidth
		}
		ratio := float64(imgWidth) / float64(width)
		rows := int(float64(imgHeight) / ratio)
		for i := 1; i < rows; i += 2 {
			fmt.Println("")
		}
		cursorUp(rows / 2)
		for i := 1; i < rows; i += 2 {
			for j := 0; j < width; j++ {
				// TODO: get average color of the area instead
				// of one pixel?
				x := int(ratio * float64(j))
				yTop := int(ratio * float64(i-1))
				yBottom := int(ratio * float64(i))
				top := ansirgb.Convert(img.At(x, yTop))
				bottom := ansirgb.Convert(img.At(x, yBottom))
				b := &block{}
				b.bottom = bottom
				// Foreground colors are lighter in some terminals.
				// Ignore top (FG) if it's the same color anyway
				if top.Code != bottom.Code {
					b.top = top
				}
				fmt.Printf("%s", b)
			}
			reset()
			fmt.Printf("\n")
		}
	}
}
