package gorevisit

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

// ImageRevisitor given a RevisitMsg and an image transformation function, runs the
// image data through the transformation and returns a new RevisitMsg with the
// transformed image
func ImageRevisitor(m *RevisitMsg, t func(src image.Image, dst image.RGBA) error) (*RevisitMsg, error) {
	reader := m.Content.ByteReader()
	srcImg, _, err := image.Decode(reader)
	if err != nil {
		return m, err
	}

	dstImg := image.NewRGBA(srcImg.Bounds())
	err = t(srcImg, *dstImg)

	dstImgBuf := bytes.NewBuffer(nil)

	format := m.ImageType()
	log.Infof("Processing image in format: %s", format)

	switch format {
	case "image/jpeg":
		err = jpeg.Encode(dstImgBuf, dstImg, nil)
		if err != nil {
			return m, err
		}
	case "image/jpg":
		err = jpeg.Encode(dstImgBuf, dstImg, nil)
		if err != nil {
			return m, err
		}

	case "image/png":
		err = png.Encode(dstImgBuf, dstImg)
		if err != nil {
			return m, err
		}
	case "image/gif":
		err = gif.Encode(dstImgBuf, dstImg, nil)
		if err != nil {
			return m, err
		}
	default:
		return m, fmt.Errorf("%s is not a supported image format", format)
	}

	dstImgBase64 := base64.StdEncoding.EncodeToString(dstImgBuf.Bytes())

	return &RevisitMsg{
		Content: ImageData{
			Data: fmt.Sprintf("data:image/%s;base64,%s", format, dstImgBase64),
		},
	}, nil
}
