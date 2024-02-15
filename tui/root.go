package tui

import (
	"ps2manager/manager"

	"github.com/rivo/tview"
)

func Menu(games []manager.GameConfig) *tview.Box {
	root := tview.NewBox().SetTitle("PS2 Game Manager").SetBorder(true)
	return root
}
