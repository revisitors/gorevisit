package gorevisit

import (
	"errors"
	"image"
	"image/gif"
)

// RevisitImage can hold either a PNG, JPEG, or GIF
type RevisitImage struct {
	images []image.Image
}

// NewRevisitImageFromMsg constructs a RevisitImage from the
// contents of a RevisitMsg
func NewRevisitImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	ri := &RevisitImage{
		images: make([]image.Image, 1),
	}

	switch r.ImageType() {
	case "image/jpeg":
		img, _, err := image.Decode(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		ri.images = append(ri.images, img)
		return ri, nil

	case "image/png":
		img, _, err := image.Decode(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		ri.images = append(ri.images, img)
		return ri, nil

	case "image/gif":
		gifs, err := gif.DecodeAll(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		for _, g := range gifs.Image {
			ri.images = append(ri.images, image.Image(g))
		}
		return ri, nil

	default:
		return nil, errors.New("invalid image type")
	}
}
