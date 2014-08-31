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
	rgbas     []image.RGBA
	palette   []color.Palette
	delay     []int
	loopCount int
	imgType   string
}

// NewRevisitImageFromMsg constructs a RevisitImage from the
// contents of a RevisitMsg
func NewRevisitImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	ri := &RevisitImage{
		rgbas:     make([]image.RGBA, 0),
		palette:   make([]color.Palette, 0),
		delay:     make([]int, 0),
		loopCount: 0,
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

		ri.rgbas = append(ri.rgbas, *rgba)
		ri.delay = append(ri.delay, 0)
		ri.loopCount = 0
		ri.imgType = r.ImageType()

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

			ri.rgbas = append(ri.rgbas, *rgba)
			ri.palette = append(ri.palette, src.Palette)
			ri.delay = append(ri.delay, 0)
		}
		ri.loopCount = gifs.LoopCount
		ri.imgType = r.ImageType()

		return ri, nil

	default:
		return nil, errors.New("invalid image type")
	}
}

func (ri *RevisitImage) RevisitMsg() (*RevisitMsg, error) {
	buf := bytes.NewBuffer(nil)
	switch ri.imgType {
	case "image/jpeg":
		err := jpeg.Encode(buf, image.Image(image.Image(&ri.rgbas[0])), nil)
		if err != nil {
			return nil, err
		}

	case "image/png":
		err := png.Encode(buf, image.Image(&ri.rgbas[0]))
		if err != nil {
			return nil, err
		}

	case "image/gif":
		g := &gif.GIF{
			Image:     make([]*image.Paletted, 0),
			LoopCount: ri.loopCount,
			Delay:     make([]int, 0),
		}

		for index, src := range ri.rgbas {
			b := src.Bounds()
			pal := image.NewPaletted(image.Rect(0, 0, b.Dx(), b.Dy()), ri.palette[index])
			draw.Draw(pal, pal.Bounds(), image.Image(&src), b.Min, draw.Src)

			g.Image = append(g.Image, pal)
			g.Delay = append(g.Delay, ri.delay[index])
		}

		buf := bytes.NewBuffer(nil)
		err := gif.EncodeAll(buf, g)
		if err != nil {
			return nil, err
		}

		dstImgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		return &RevisitMsg{
			Content: ImageData{
				Data: fmt.Sprintf("data:%s;base64,%s", ri.imgType, dstImgBase64),
			},
		}, nil

	default:
		return nil, errors.New("invalid image type")
	}

	dstImgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return &RevisitMsg{
		Content: ImageData{
			Data: fmt.Sprintf("data:%s;base64,%s", ri.imgType, dstImgBase64),
		},
	}, nil
}
