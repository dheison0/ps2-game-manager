package tui

import (
	"ps2manager/manager"

	"github.com/rivo/tview"
)

func SetupMenu() {
	menu = tview.NewList()
	menu.SetBorder(true)
	menu.SetBorderPadding(1, 1, 1, 1)
	menu.SetTitle(" PS2 Game Manager ")
	menu.SetTitleAlign(tview.AlignCenter)
	menu.ShowSecondaryText(false)

	games := manager.GetAll()
	menu.Clear()
	for _, g := range games {
		AddMenuItem(g)
	}
	menu.AddItem("Install", "", 'i', func() {
		SelectFile(func(f string) {
			UpdateInstallForm(f)
			pages.SwitchToPage("installForm")
		})
	})
	menu.AddItem("Get covers", "", 'c', func() {})
	menu.AddItem("Quit", "", 'q', func() { app.Stop() })
}

func RemoveMenuItem(index int) {
	menu.RemoveItem(index)
}

func AddMenuItem(game *manager.GameConfig) {
	menu.AddItem(game.GetName(), "", 0, nil)
}
