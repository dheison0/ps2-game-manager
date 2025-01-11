package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"math"

	"github.com/nfnt/resize"
)

// ResizeJPGKeepingAspectRatio do what the name says :)
// it reads an JPG image file and resize it keeping original aspect ratio and
// taking advantage of as much resolution as possible
func ResizeJPGKeepingAspectRatio(data []byte, maxWidth, maxHeight int) ([]byte, error) {
	empty := []byte{}
	cover, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return empty, err
	}
	if cover.Width > maxWidth || cover.Height > maxHeight {
		mw := float64(maxWidth)
		mh := float64(maxHeight)
		aspect := float64(cover.Width) / float64(cover.Height)
		// get height based on width
		width := mw
		height := math.Floor(float64(mw) / aspect)
		if height > mh {
			// if height obtained is too big, make it the max height and then get
			// width based on height
			width = math.Floor(aspect * mh)
			height = mh
		}
		data, err = ResizeJPG(bytes.NewReader(data), uint(width), uint(height))
		if err != nil {
			return empty, err
		}
	}
	return data, nil
}

// ResizeJPG resizes the image to the defined width and height
func ResizeJPG(data io.Reader, width, height uint) ([]byte, error) {
	var result []byte
	originalImage, err := jpeg.Decode(data)
	if err != nil {
		return result, err
	}
	newImage := resize.Resize(width, height, originalImage, resize.Lanczos3)
	buffer := bytes.NewBuffer(result)
	err = jpeg.Encode(buffer, newImage, nil)
	if err != nil {
		return result, err
	}
	return buffer.Bytes(), nil
}
