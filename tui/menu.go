package tui

import (
	"ps2manager/manager"

	"github.com/rivo/tview"
)

type MenuScreen struct {
	root *tview.List
}

func NewMenuScreen() *MenuScreen {
	screen := &MenuScreen{}
	screen.root = tview.NewList()
	screen.root.
		ShowSecondaryText(false).
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle(" PS2 Game Manager ").
		SetTitleAlign(tview.AlignCenter)
	screen.root.SetSelectedFunc(func(index int, _, _ string, s rune) {
		if s != 0 {
			return
		}
		actionsMenu.UpdateGame(index)
		actionsMenu.Show()
	})
	screen.UpdateItemList()
	pages.AddPage("menu", screen.root, true, false)
	return screen
}

func (m *MenuScreen) Show() {
	pages.SwitchToPage("menu")
}

func (m *MenuScreen) RemoveItem(index int) {
	m.root.RemoveItem(index)
}

func (m *MenuScreen) AddItem(game *manager.GameConfig) {
	m.root.AddItem(game.GetName(), "", 0, nil)
}

func (m *MenuScreen) UpdateItemList() {
	games := manager.GetAll()
	m.root.Clear()
	for _, g := range games {
		m.AddItem(g)
	}
	m.root.AddItem("Install", "", 'i', func() {
		fileSelector.SetSelectFileFunc(func(f string) {
			install.NewForm(f)
			install.Show()
		})
		fileSelector.Show()
	})
	m.root.
		AddItem("Get covers", "", 'c', func() {
			covers.Show()
			covers.DownloadMissingCovers()
		}).
		AddItem("Quit", "", 'q', func() { app.Stop() })
}
