package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/hooklift/iso9660"
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

func ReadFileFromISO(iso, filename string) ([]byte, error) {
	isoFile, err := os.Open(iso)
	if err != nil {
		fmt.Println("Não foi possivel abrir a iso")
		return []byte{}, err
	}
	isoReader, err := iso9660.NewReader(isoFile)
	if err != nil {
		return []byte{}, err
	}
	var wantedFile fs.FileInfo
	for {
		f, err := isoReader.Next()
		if err == io.EOF {
			return []byte{}, errors.New("file not found")
		} else if err != nil {
			return []byte{}, err
		} else if strings.ToLower(f.Name()) == strings.ToLower(filename) {
			wantedFile = f
			break
		}
	}
	fReader := wantedFile.Sys().(io.Reader)
	return io.ReadAll(fReader)
}
