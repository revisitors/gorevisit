package gorevisit

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
)

// RevisitImage can hold either a PNG, JPEG, or GIF
type RevisitImage struct {
	Rgbas     []image.RGBA
	Palette   []color.Palette
	Delay     []int
	LoopCount int
	ImgType   string
}

// NewRevisitImageFromMsg constructs a RevisitImage from the
// contents of a RevisitMsg
func NewRevisitImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	ri := &RevisitImage{
		Rgbas:     make([]image.RGBA, 0),
		Palette:   make([]color.Palette, 0),
		Delay:     make([]int, 0),
		LoopCount: 0,
	}

	switch r.ImageType() {
	case "image/jpeg", "image/png":
		src, _, err := image.Decode(r.ImageByteReader())
		if err != nil {
			return nil, err
		}

		b := src.Bounds()
		rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), src, b.Min, draw.Src)

		ri.Rgbas = append(ri.Rgbas, *rgba)
		ri.Delay = append(ri.Delay, 0)
		ri.LoopCount = 0
		ri.ImgType = r.ImageType()

		return ri, nil

	case "image/gif":
		gifs, err := gif.DecodeAll(r.ImageByteReader())
		if err != nil {
			return nil, err
		}

		for _, src := range gifs.Image {
			b := src.Bounds()
			rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(rgba, rgba.Bounds(), src, b.Min, draw.Src)

			ri.Rgbas = append(ri.Rgbas, *rgba)
			ri.Palette = append(ri.Palette, src.Palette)
			ri.Delay = append(ri.Delay, 0)
		}
		ri.LoopCount = gifs.LoopCount
		ri.ImgType = r.ImageType()

		return ri, nil

	default:
		return nil, errors.New("invalid image type")
	}
}

func (ri *RevisitImage) Transform(t func(src draw.Image)) {
	for _, frame := range ri.Rgbas {
		t(draw.Image(&frame))
	}
}

func (ri *RevisitImage) RevisitMsg() (*RevisitMsg, error) {
	buf := bytes.NewBuffer(nil)

	switch ri.ImgType {
	case "image/jpeg":
		err := jpeg.Encode(buf, image.Image(image.Image(&ri.Rgbas[0])), nil)
		if err != nil {
			return nil, err
		}

	case "image/png":
		err := png.Encode(buf, image.Image(&ri.Rgbas[0]))
		if err != nil {
			return nil, err
		}

	case "image/gif":
		g := &gif.GIF{
			Image:     make([]*image.Paletted, 0),
			LoopCount: ri.LoopCount,
			Delay:     make([]int, 0),
		}

		for index, src := range ri.Rgbas {
			b := src.Bounds()
			pal := image.NewPaletted(image.Rect(0, 0, b.Dx(), b.Dy()), ri.Palette[index])
			draw.Draw(pal, pal.Bounds(), image.Image(&src), b.Min, draw.Src)

			g.Image = append(g.Image, pal)
			g.Delay = append(g.Delay, ri.Delay[index])
		}

		buf := bytes.NewBuffer(nil)
		err := gif.EncodeAll(buf, g)
		if err != nil {
			return nil, err
		}

		dstImgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		return &RevisitMsg{
			Content: ImageData{
				Data: fmt.Sprintf("data:%s;base64,%s", ri.ImgType, dstImgBase64),
			},
		}, nil

	default:
		return nil, errors.New("invalid image type")
	}

	dstImgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return &RevisitMsg{
		Content: ImageData{
			Data: fmt.Sprintf("data:%s;base64,%s", ri.ImgType, dstImgBase64),
		},
	}, nil
}
