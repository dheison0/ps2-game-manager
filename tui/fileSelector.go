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
		fi := fileSelector.GetCurrentItem()
		if fi == 0 {
			fileSelectorActualPath = path.Join(fileSelectorActualPath, "..")
		} else {
			selected := folders[fi-1].Name() // -1 because first is up dir
			fileSelectorActualPath = path.Join(fileSelectorActualPath, selected)
		}
		UpdateFileList()
	}
	done := func() {
		// file index is current - length of folders - 1 (of up dir)
		fi := fileSelector.GetCurrentItem() - len(folders) - 1
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
	fileSelector.AddItem("../", "", 0, chdir)
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
