package utils

import (
	"bytes"
	"image/jpeg"
	"io"

	"github.com/nfnt/resize"
)

func BytesToString(data []byte) string {
	n := bytes.IndexByte(data, 0)
	if n < 0 {
		n = len(data) - 1
	}
	return string(data[:n])
}

func ResizeJPG(data io.Reader, width, height uint) ([]byte, error) {
	var result []byte
	oldImage, err := jpeg.Decode(data)
	if err != nil {
		return result, err
	}
	newImage := resize.Resize(width, height, oldImage, resize.Lanczos3)
	buffer := bytes.NewBuffer(result)
	err = jpeg.Encode(buffer, newImage, nil)
	if err != nil {
		return result, err
	}
	return buffer.Bytes(), nil
}
