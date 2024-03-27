package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/rivo/tview"
)

type ActionsMenuScreen struct {
	root      *tview.Flex
	info      *tview.TextView
	buttons   *tview.List
	gameIndex int
}

func NewActionsMenuScreen() *ActionsMenuScreen {
	screen := &ActionsMenuScreen{}
	screen.root = tview.NewFlex()
	screen.root.
		SetDirection(tview.FlexRow).
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2).
		SetTitle(" Game Actions ").
		SetTitleAlign(tview.AlignCenter)
	screen.info = tview.NewTextView()
	screen.info.
		SetDynamicColors(true).
		SetBorderPadding(0, 1, 0, 0)
	screen.buttons = tview.NewList()
	screen.buttons.
		ShowSecondaryText(false).
		AddItem("Rename", "", 'r', func() {}).
		AddItem("Delete", "", 'd', func() {
			actionsDelete.SetGameIndex(screen.gameIndex)
			actionsDelete.Show()
		}).
		AddItem("Go back", "", 'q', menu.Show)
	screen.root.
		AddItem(screen.info, 3, 1, false).
		AddItem(screen.buttons, 0, 1, true)
	pages.AddPage("actionsMenu", screen.root, true, false)
	return screen
}

func (m *ActionsMenuScreen) Show() {
	pages.SwitchToPage("actionsMenu")
}

func (m *ActionsMenuScreen) UpdateGame(gameIndex int) {
	if gameIndex != -1 { // keep the same if it is -1
		m.gameIndex = gameIndex
	}
	game := manager.Get(m.gameIndex)
	m.info.SetText(fmt.Sprintf(
		"[white]Name: [purple]%s\n[white]Image: [purple]%s[default]",
		game.GetName(),
		game.GetImage(),
	))
}

type ActionsDeleteScreen struct {
	root      *tview.Modal
	gameIndex int
}

func NewActionsDeleteScreen() *ActionsDeleteScreen {
	screen := &ActionsDeleteScreen{}
	screen.root = tview.NewModal()
	screen.root.
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2).
		SetTitle(" Delete game ").
		SetTitleAlign(tview.AlignCenter)
	screen.root.AddButtons([]string{"no", "yes"})
	screen.root.SetDoneFunc(func(_ int, name string) {
		if name == "yes" {
			if err := manager.Delete(screen.gameIndex); err != nil {
				errorDialog.SetMessage("Failed to delete game: " + err.Error())
				errorDialog.Show()
				return
			}
			menu.RemoveItem(screen.gameIndex)
			menu.Show()
			return
		}
		actionsMenu.Show()
	})
	pages.AddPage("actionsDelete", screen.root, true, false)
	return screen
}

func (d *ActionsDeleteScreen) Show() {
	pages.SwitchToPage("actionsDelete")
}

func (d *ActionsDeleteScreen) SetGameIndex(gameIndex int) {
	d.gameIndex = gameIndex
	game := manager.Get(d.gameIndex)
	d.root.SetText("Do you really want to delete " + game.GetName() + "?")
	d.root.SetFocus(0)
}
