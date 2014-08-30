package gorevisit

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

// RevisitImage can hold either a PNG, JPEG, or GIF
type RevisitImage struct {
	PNG  *image.Image
	JPEG *image.Image
	GIF  *gif.GIF
	Type string
}

// IsGif return true if the RevisitImage holds a GIF
func (ri *RevisitImage) IsGIF() bool {
	switch ri.Type {
	case "image/gif":
		return true
	default:
		return false
	}
}

// IsJPEG return true if the RevisitImage holds a JPEG,
func (ri *RevisitImage) IsJPEG() bool {
	switch ri.Type {
	case "image/jpeg":
		return true
	default:
		return false
	}
}

// IsJPEG return true if the RevisitImage holds a PNG,
func (ri *RevisitImage) IsPNG() bool {
	switch ri.Type {
	case "image/png":
		return true
	default:
		return false
	}
}

// NewRevisitImageFromMsg constructs a RevisitImage from the
// contents of a RevisitMsg
func NewRevisitImageFromMsg(r *RevisitMsg) (*RevisitImage, error) {
	switch r.ImageType() {
	case "image/gif":
		img, err := gif.DecodeAll(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		return &RevisitImage{
			GIF:  img,
			Type: r.ImageType(),
		}, nil

	case "image/jpeg":
		img, err := jpeg.Decode(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		return &RevisitImage{
			JPEG: &img,
			Type: r.ImageType(),
		}, nil

	case "image/png":
		img, err := png.Decode(r.Content.ByteReader())
		if err != nil {
			return nil, err
		}

		return &RevisitImage{
			PNG:  &img,
			Type: r.ImageType(),
		}, nil

	default:
		return nil, errors.New("invalid image type")
	}
}
