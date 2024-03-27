package tui

import (
	"github.com/rivo/tview"
)

type ErrorDialogScreen struct {
	root *tview.Modal
}

func NewErrorDialogScreen() *ErrorDialogScreen {
	screen := &ErrorDialogScreen{}
	screen.root = tview.NewModal()
	screen.root.
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2).
		SetTitleAlign(tview.AlignCenter).
		SetTitle(" !!! ERROR !!!")
	screen.root.
		AddButtons([]string{"done"}).
		SetDoneFunc(func(_ int, _ string) { menu.Show() })
	pages.AddPage("errorDialog", screen.root, true, false)
	return screen
}

func (e *ErrorDialogScreen) Show() {
	pages.SwitchToPage("errorDialog")
}

func (e *ErrorDialogScreen) SetMessage(message string) {
	e.root.SetText(message)
}
