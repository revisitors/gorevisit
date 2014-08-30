package gorevisit

import (
	"errors"
	"image"
	"image/gif"
)

// RevisitImage can hold either a PNG, JPEG, or GIF
type RevisitImage struct {
	images    []image.Image
	palettes  []image.PalettedImage
	delay     []int
	loopCount int
}

// NewRevisitImageFromMsg constructs a RevisitImage from the
// contents of a RevisitMsg
func NewRevisitImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	ri := &RevisitImage{
		images: make([]image.Image, 0),
	}

	switch r.ImageType() {
	case "image/jpeg":
		img, _, err := image.Decode(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		ri.images = append(ri.images, img)
		ri.delay = append(ri.delay, 0)
		ri.loopCount = 0

		return ri, nil

	case "image/png":
		img, _, err := image.Decode(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		ri.images = append(ri.images, img)
		ri.delay = append(ri.delay, 0)
		ri.loopCount = 0

		return ri, nil

	case "image/gif":
		gifs, err := gif.DecodeAll(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		for i, g := range gifs.Image {
			ri.palettes = append(ri.palettes, g)
			ri.delay = append(ri.delay, gifs.Delay[i])
		}
		ri.loopCount = gifs.LoopCount
		return ri, nil

	default:
		return nil, errors.New("invalid image type")
	}
}
