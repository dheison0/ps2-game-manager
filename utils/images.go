package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	"github.com/nfnt/resize"
)

func ResizeJPGToMax(data []byte, width, height int) ([]byte, error) {
	empty := []byte{}
	cover, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return empty, err
	}
	if cover.Width > width || cover.Height > height {
		data, err = ResizeJPG(bytes.NewReader(data), uint(width), uint(height))
		if err != nil {
			return empty, err
		}
	}
	return data, nil
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
