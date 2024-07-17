package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const (
	KIBI_BYTE = 1_024
	MEBI_BYTE = 1_048_576
	GIBI_BYTE = 1_073_741_824
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// SortDirItems sort files and folders by putting folders before files and sorting by name
func SortDirItems(first, second fs.DirEntry) int {
	firstName := strings.ToLower(first.Name())
	secondName := strings.ToLower(second.Name())
	// if entry "first" is a directory and "second" isn't or "first" name is
	// lower than "second" name, "first" must be before "second"
	if first.IsDir() && !second.IsDir() || firstName < secondName {
		return -1
	}
	return 1
}

func GetFileSizeSum(files []string) (int64, error) {
	var size int64 = 0
	for _, f := range files {
		file, err := os.Stat(f)
		if err != nil {
			return 0, errors.New("error on file '" + f + "': " + err.Error())
		}
		size += file.Size()
	}
	return size, nil
}

func FileSizeToHumanReadable(size int64) string {
	s := float64(size)
	if s >= GIBI_BYTE {
		return fmt.Sprintf("%0.2fGiB", s/GIBI_BYTE)
	} else if s >= MEBI_BYTE {
		return fmt.Sprintf("%0.2fMiB", s/MEBI_BYTE)
	} else if s >= KIBI_BYTE {
		return fmt.Sprintf("%0.2fKiB", s/KIBI_BYTE)
	}
	return fmt.Sprintf("%dB", size)
}
