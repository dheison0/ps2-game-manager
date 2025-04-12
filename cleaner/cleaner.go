package cleaner

import (
	"fmt"
	"os"
	"path"
	"ps2manager/manager"
	"slices"
	"strings"
)

type GameFile struct {
	File string
	Size int64
}

func getAllGameFiles() ([]GameFile, error) {
	files := []GameFile{}
	root := manager.GetDataDir()
	dirItems, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("failed to read root directory items")
	}
	for _, i := range dirItems {
		name := i.Name()
		if name == "ul.cfg" || i.IsDir() || !strings.Contains(name, "ul.") {
			continue
		}

		filePath := path.Join(root, name)
		fileInfo, err := i.Info()

		if err != nil {
			return nil, fmt.Errorf("can't read '%s' info", filePath)
		}
		files = append(files, GameFile{File: filePath, Size: fileInfo.Size()})
	}
	return files, nil
}

func GetUnusedGameFiles() ([]GameFile, error) {
	files, err := getAllGameFiles()
	if err != nil {
		return nil, err
	}
	for _, game := range manager.GetAll() {
		for _, file := range game.Files {
			idx := slices.IndexFunc(files, func(i GameFile) bool { return i.File == file })
			if idx != -1 {
				files = slices.Delete(files, idx, idx+1)
			}
		}
	}
	return files, err
}

func DeleteFiles(files []GameFile) error {
	for _, f := range files {
		if err := os.Remove(f.File); err != nil {
			return err
		}
	}
	return nil
}
