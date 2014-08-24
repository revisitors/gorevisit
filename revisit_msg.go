package gorevisit

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"strings"
)

type Image struct {
	Data string `json:"data"`
}

func (i Image) ByteReader() io.Reader {
	dataUri := i.Data
	data := strings.Split(dataUri, ",")[1]
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
}

type Audio struct {
	Data string `json:data"`
}

type Meta struct {
	Audio Audio `json:"audio"`
}

type RevisitMsg struct {
	Content Image `json:"content"`
	Meta    Meta  `json:"meta"`
}

func ImageRevisitor(m *RevisitMsg, t func(src image.Image, dst image.RGBA) error) (*RevisitMsg, error) {
	reader := m.Content.ByteReader()
	srcImg, format, err := image.Decode(reader)
	if err != nil {
		return m, err
	}

	if format != "jpeg" {
		return m, nil
	}

	dstImg := image.NewRGBA(srcImg.Bounds())
	err = t(srcImg, *dstImg)

	dstImgBuf := bytes.NewBuffer(nil)
	err = jpeg.Encode(dstImgBuf, dstImg, nil)
	if err != nil {
		return m, err
	}

	dstImgBase64 := base64.StdEncoding.EncodeToString(dstImgBuf.Bytes())

	return &RevisitMsg{
		Content: Image{
			Data: fmt.Sprintf("data:image/%s;base64,%s", format, dstImgBase64),
		},
	}, nil
}

func BytesToDataURI(data []byte, contentType string) string {
	return fmt.Sprintf("data:%s;base64,%s",
		contentType, base64.StdEncoding.EncodeToString(data))
}

func NewRevisitMsgFromFiles(mediaPath ...string) (*RevisitMsg, error) {
	if len(mediaPath) < 1 || len(mediaPath) > 2 {
		return &RevisitMsg{}, errors.New("must have image, may have audio")
	}

	imageBytes, err := ioutil.ReadFile(mediaPath[0])
	if err != nil {
		return &RevisitMsg{}, err
	}

	// FIXME: add image type detection instead of hard coded jpeg
	imageDataURI := BytesToDataURI(imageBytes, "image/jpeg")

	soundBytes, _ := ioutil.ReadFile(mediaPath[1])
	if err != nil {
		return &RevisitMsg{}, err
	}
	// FIXME: add sound type detection instead of hard coded ogg
	soundDataURI := BytesToDataURI(soundBytes, "audio/ogg")

	content := &Image{
		Data: imageDataURI,
	}

	audioContent := &Audio{
		Data: soundDataURI,
	}

	metaContent := &Meta{
		Audio: *audioContent,
	}

	revisitMsg := &RevisitMsg{
		Content: *content,
		Meta:    *metaContent,
	}
	return revisitMsg, nil
}
