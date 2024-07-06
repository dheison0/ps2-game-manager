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
	callback func(string)
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
	var items []os.DirEntry
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue // skip hidden files
		}
		items = append(items, entry)
	}
	slices.SortFunc(items, utils.SortDirItems)
	f.root.
		Clear().
		AddItem("../", "", 0, nil)
	for _, item := range items {
		name := item.Name()
		if item.IsDir() {
			name += "/"
		}
		f.root.AddItem(name, "", 0, nil)
	}
	f.root.AddItem("Quit!", "", 'q', nil)
}

// SetSelectFileFunc is a function that set a callback function called when a file is selected
func (f *FileSelectorScreen) SetSelectFileFunc(selectedFunc func(string)) {
	f.callback = selectedFunc
}
