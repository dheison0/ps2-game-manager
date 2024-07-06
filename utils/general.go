package utils

import (
	"io/fs"
	"os"
	"strings"
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
