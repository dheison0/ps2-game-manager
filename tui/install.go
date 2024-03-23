package tui

import (
	"ps2manager/manager"

	"github.com/rivo/tview"
)

var gameName, isoFilePath string

func SetupInstall() {
	installForm = tview.NewForm()
	installForm.SetBorder(true)
	installForm.SetBorderPadding(1, 1, 1, 1)
	installForm.SetTitle(" Install new game ")
	installForm.AddButton("Install", func() {})
	installForm.AddButton("Cancel", func() { pages.SwitchToPage("menu") })
}

func UpdateInstallForm(isoFile string) {
	isoFilePath = isoFile
	gameName = ""
	installForm.Clear(false)
	installForm.AddTextView("ISO file:", isoFile, 0, 1, false, false)
	installForm.AddInputField(
		"Game Name:",
		gameName,
		manager.MaxNameSize,
		func(t string, _ rune) bool { return len(t) <= manager.MaxNameSize },
		func(t string) { gameName = t },
	)
}
