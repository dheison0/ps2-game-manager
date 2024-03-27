package tui

import (
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/rivo/tview"
)

type FileSelectorScreen struct {
	root       *tview.List
	callback   func(string)
	actualPath string
}

func NewFileSelectorScreen() *FileSelectorScreen {
	screen := &FileSelectorScreen{}
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
	pages.SwitchToPage("fileSelector")
}

func (f *FileSelectorScreen) UpdateFileList() {
	entries, err := os.ReadDir(f.actualPath)
	if err != nil {
		return
	}
	var items []os.DirEntry
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue // skip hidden files
		}
		items = append(items, entry)
	}
	slices.SortFunc(items, func(a, b fs.DirEntry) int { // Put dirs at top and sort by name
		// if 'a' is a dir and 'b' isn't then 'a' must be first or
		// if 'a' is lower than 'b' it must be least
		if a.IsDir() && !b.IsDir() || a.Name() < b.Name() {
			return -1
		}
		return 1
	})
	f.root.
		Clear().
		AddItem("../", "", 0, nil)
	for _, i := range items {
		name := i.Name()
		if i.IsDir() {
			name += "/"
		}
		f.root.AddItem(name, "", 0, nil)
	}
	f.root.AddItem("Quit!", "", 'q', nil)
}

func (f *FileSelectorScreen) SetSelectFileFunc(selectedFunc func(string)) {
	f.callback = selectedFunc
}
