package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ErrorDialog(message string) {
	dialog := tview.NewModal()
	dialog.SetTitle("ERROR!!!")
	dialog.SetBackgroundColor(tcell.Color100)
	dialog.SetText(message)
	dialog.AddButtons([]string{"Back"})
	dialog.SetDoneFunc(func(_ int, _ string) { Menu() })
	pages.AddAndSwitchToPage("Error", dialog, true)
}
