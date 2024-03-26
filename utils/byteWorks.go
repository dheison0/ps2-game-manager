package utils

// !!! This file must have only functions that works directly with bytes !!!

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/hooklift/iso9660"
)

var (
	ErrUnableToOpen = errors.New("unable to open file")
	ErrNotFound     = errors.New("file not found")
)

func BytesToString(data []byte) string {
	n := bytes.IndexByte(data, 0)
	if n < 0 {
		n = len(data) - 1
	}
	return string(data[:n])
}

func ReadFileFromISO(iso, filename string) ([]byte, error) {
	empty := []byte{}
	isoFile, err := os.Open(iso)
	if err != nil {
		return empty, ErrUnableToOpen
	}
	isoReader, err := iso9660.NewReader(isoFile)
	if err != nil {
		return empty, err
	}
	var wantedFile fs.FileInfo
	for {
		f, err := isoReader.Next()
		if err == io.EOF {
			return empty, ErrNotFound
		} else if err != nil {
			return empty, err
		} else if strings.EqualFold(f.Name(), filename) {
			wantedFile = f
			break
		}
	}
	fReader := wantedFile.Sys().(io.Reader)
	return io.ReadAll(fReader)
}
