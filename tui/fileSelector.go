package tui

import (
	"os"
	"path"
	"ps2manager/utils"
	"slices"
	"strings"

	"github.com/rivo/tview"
)

type FileSelectorScreen struct {
	screenId   string
	root       *tview.List
	actualPath string

	// callback is a function called when a file is selected
	callback       func(string)
	acceptFolders  bool
	fileExtensions string
}

func NewFileSelectorScreen() *FileSelectorScreen {
	screen := &FileSelectorScreen{screenId: "fileSelector"}
	screen.root = tview.NewList()
	root, err := os.Getwd()
	if err == nil {
		screen.actualPath = root
	}
	screen.root.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle(" Select file ")
	screen.root.
		ShowSecondaryText(false).
		SetSelectedFunc(func(_ int, name, _ string, shortcut rune) {
			if shortcut == 'q' {
				menu.Show()
				return
			}
			selectedItemFullPath := path.Join(screen.actualPath, name)
			if strings.Contains(name, "/") {
				screen.actualPath = selectedItemFullPath
				screen.UpdateFileList()
			} else if name == "." && screen.acceptFolders {
				screen.callback(screen.actualPath)
			} else {
				screen.callback(selectedItemFullPath)
			}
		})
	screen.UpdateFileList()
	pages.AddPage("fileSelector", screen.root, true, false)
	return screen
}

func (f *FileSelectorScreen) Show() {
	pages.SwitchToPage(f.screenId)
}

func (f *FileSelectorScreen) UpdateFileList() {
	entries, err := os.ReadDir(f.actualPath)
	if err != nil {
		f.actualPath = path.Join(f.actualPath, "..")
		return
	}
	items := filterDirItems(entries, f.fileExtensions, false)
	slices.SortFunc(items, utils.SortDirItems)
	f.root.Clear()
	if f.acceptFolders {
		f.root.AddItem(".", "", 0, nil)
	}
	f.root.AddItem("../", "", 0, nil)
	for _, item := range items {
		name := item.Name()
		if item.IsDir() {
			name += "/"
		}
		f.root.AddItem(name, "", 0, nil)
	}
	f.root.AddItem("Quit!", "", 'q', nil)
}

func filterDirItems[T []os.DirEntry](entries T, fileExtensions string, showHiddenFiles bool) T {
	var items []os.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		if (!showHiddenFiles && isHiddenFile(name)) || (!entry.IsDir() && fileExtensions != "" && !fileExtensionInList(fileExtensions, name)) {
			continue // skip hidden files
		}
		items = append(items, entry)
	}
	return items
}

func isHiddenFile(name string) bool {
	return strings.HasPrefix(name, ".")
}

func fileExtensionInList(fileExtensions, fileName string) bool {
	nameParts := strings.Split(fileName, ".")
	extension := strings.ToLower(nameParts[len(nameParts)-1])
	return strings.Contains(fileExtensions, extension)
}

// SetSelectFileFunc is a function that set a callback function called when a file is selected
func (f *FileSelectorScreen) SetSelectFileConfig(selectedFunc func(string), acceptFolders bool, fileExtensions string) {
	f.callback = selectedFunc
	f.acceptFolders = acceptFolders
	f.fileExtensions = strings.ToLower(fileExtensions)
	f.UpdateFileList()
}
