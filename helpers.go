package gorevisit

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// BytesToDataURI returns a data URI encoded string given a byte array and a content type
// See RFC2397 - http://tools.ietf.org/html/rfc2397
func BytesToDataURI(data []byte, contentType string) string {
	return fmt.Sprintf("data:%s;base64,%s",
		contentType, base64.StdEncoding.EncodeToString(data))
}

// DataURIToDecodedContent returns a content type string and an array of bytes
// given a data URI encoded string.
// See RFC2397 - http://tools.ietf.org/html/rfc2397
func DataURIToDecodedContent(dataURI string) (*DecodedContent, error) {
	parts := strings.Split(dataURI, ",")
	contentType := parts[0]
	contentBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return &DecodedContent{}, err
	}
	return &DecodedContent{Type: contentType, Data: contentBytes}, nil
}

func getFormat(file *os.File) string {
	bytes := make([]byte, 4)
	n, _ := file.ReadAt(bytes, 0)
	if n < 4 {
		return ""
	}
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return "image/png"
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return "image/jpg"
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return "image/gif"
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return "image/bmp"
	}
	return ""
}

func getGifDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 4)
	file.ReadAt(bytes, 6)
	width = int(bytes[0]) + int(bytes[1])*256
	height = int(bytes[2]) + int(bytes[3])*256
	return
}

func getBmpDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 18)
	width = int(bytes[3])<<24 | int(bytes[2])<<16 | int(bytes[1])<<8 | int(bytes[0])
	height = int(bytes[7])<<24 | int(bytes[6])<<16 | int(bytes[5])<<8 | int(bytes[4])
	return
}

func getPngDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 16)
	width = int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	height = int(bytes[4])<<24 | int(bytes[5])<<16 | int(bytes[6])<<8 | int(bytes[7])
	return
}

func getJpgDimensions(file *os.File) (width int, height int) {
	fi, _ := file.Stat()
	fileSize := fi.Size()

	position := int64(4)
	bytes := make([]byte, 4)
	file.ReadAt(bytes[:2], position)
	length := int(bytes[0]<<8) + int(bytes[1])
	for position < fileSize {
		position += int64(length)
		file.ReadAt(bytes, position)
		length = int(bytes[2])<<8 + int(bytes[3])
		if (bytes[1] == 0xC0 || bytes[1] == 0xC2) && bytes[0] == 0xFF && length > 7 {
			file.ReadAt(bytes, position+5)
			width = int(bytes[2])<<8 + int(bytes[3])
			height = int(bytes[0])<<8 + int(bytes[1])
			return
		}
		position += 2
	}
	return 0, 0
}
