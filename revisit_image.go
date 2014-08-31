package gorevisit

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

// RevisitImage can hold either a PNG, JPEG, or GIF
type RevisitImage struct {
	images    []image.Image
	palettes  []image.Paletted
	delay     []int
	loopCount int
	imgType   string
}

func newImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	ri := &RevisitImage{
		images:    make([]image.Image, 0),
		palettes:  make([]image.Paletted, 0),
		delay:     make([]int, 0),
		loopCount: 0,
	}

	img, _, err := image.Decode(r.ImageByteReader())
	if err != nil {
		return nil, err
	}

	ri.images = append(ri.images, img)
	ri.delay = append(ri.delay, 0)
	ri.loopCount = 0
	ri.imgType = r.ImageType()

	return ri, nil
}

func newPalettedImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	ri := &RevisitImage{
		images: make([]image.Image, 0),
	}

	gifs, err := gif.DecodeAll(r.ImageByteReader())
	if err != nil {
		return nil, err
	}

	for i, g := range gifs.Image {
		ri.palettes = append(ri.palettes, *g)
		ri.delay = append(ri.delay, gifs.Delay[i])
	}
	ri.loopCount = gifs.LoopCount
	ri.imgType = r.ImageType()

	return ri, nil

}

// NewRevisitImageFromMsg constructs a RevisitImage from the
// contents of a RevisitMsg
func NewRevisitImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	switch r.ImageType() {
	case "image/jpeg":
		return newImageFromMsg(r)

	case "image/png":
		return newImageFromMsg(r)

	case "image/gif":
		return newPalettedImageFromMsg(r)

	default:
		return nil, errors.New("invalid image type")
	}
}

func (ri *RevisitImage) palettedToRevisitMsg() (*RevisitMsg, error) {
	g := &gif.GIF{
		Image:     make([]*image.Paletted, 0),
		LoopCount: ri.loopCount,
		Delay:     make([]int, 0),
	}

	for index, pal := range ri.palettes {
		g.Image = append(g.Image, &pal)
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
}

func (ri *RevisitImage) imageToRevisitMsg() (*RevisitMsg, error) {
	buf := bytes.NewBuffer(nil)
	switch ri.imgType {
	case "image/jpeg":
		err := jpeg.Encode(buf, ri.images[0], nil)
		if err != nil {
			return nil, err
		}

	case "image/png":
		err := png.Encode(buf, ri.images[0])
		if err != nil {
			return nil, err
		}

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

// RevisitMsg returns a RevisitMsg from a RevisitImage
func (ri *RevisitImage) RevisitMsg() (*RevisitMsg, error) {
	switch ri.imgType {
	case "image/gif":
		return ri.palettedToRevisitMsg()
	default:
		return ri.imageToRevisitMsg()
	}
}
