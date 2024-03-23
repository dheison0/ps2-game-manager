package tui

import (
	"os"
	"path"
	"strings"

	"github.com/rivo/tview"
)

var fileSelectorActualPath = "."
var callback func(string)

func SetupFileSelector() {
	root, err := os.Getwd()
	if err == nil {
		fileSelectorActualPath = root
	}
	fileSelector = tview.NewList()
	fileSelector.SetBorder(true)
	fileSelector.SetBorderPadding(1, 1, 1, 1)
	fileSelector.SetTitle(" Select file ")
	fileSelector.ShowSecondaryText(false)
	UpdateFileList()
}

func UpdateFileList() {
	entries, err := os.ReadDir(fileSelectorActualPath)
	if err != nil {
		return
	}
	var files, folders []os.DirEntry
	chdir := func() {
		fi := fileSelector.GetCurrentItem() - 1 // -1 here is necessary because it has ".." as the first folder
		selected := folders[fi].Name()
		fileSelectorActualPath = path.Join(fileSelectorActualPath, selected)
		UpdateFileList()
	}
	done := func() {
		fi := fileSelector.GetCurrentItem() - len(folders) - 1 // -1 here is necessary because it has ".." as the first folder
		selected := files[fi].Name()
		callback(path.Join(fileSelectorActualPath, selected))
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue // skip hidden files
		}
		if entry.IsDir() {
			folders = append(folders, entry)
		} else {
			files = append(files, entry)
		}
	}
	fileSelector.Clear()
	fileSelector.AddItem("../", "", 0, func() {
		fileSelectorActualPath = path.Join(fileSelectorActualPath, "..")
		UpdateFileList()
	})
	for _, d := range folders {
		fileSelector.AddItem(d.Name()+"/", "", 0, chdir)
	}
	for _, f := range files {
		fileSelector.AddItem(f.Name(), "", 0, done)
	}
	fileSelector.AddItem("Quit!", "", 'q', func() { pages.SwitchToPage("menu") })
}

func SelectFile(selectedFunc func(string)) {
	callback = selectedFunc
	pages.SwitchToPage("fileSelector")
}
