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
		AddItem("Rename", "", 'r', func() {
			actionsRename.SetGameIndex(screen.gameIndex)
			actionsRename.Show()
		}).
		AddItem("Export as ISO", "", 'e', func() {
			game := manager.Get(screen.gameIndex)
			gameExport.SetGame(game)
			gameExport.Show()
		}).
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

type ActionsRenameScreen struct {
	root      *tview.Form
	gameIndex int
	newName   string
}

func NewActionsRenameScreen() *ActionsRenameScreen {
	screen := &ActionsRenameScreen{}
	screen.root = tview.NewForm()
	screen.root.
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2).
		SetTitle(" Game Rename ").
		SetTitleAlign(tview.AlignCenter)
	screen.root.AddButton("Rename", func() {
		if manager.Get(screen.gameIndex).GetName() == screen.newName {
			actionsMenu.Show()
			return
		} else if err := manager.Rename(screen.gameIndex, screen.newName); err != nil {
			errorDialog.SetMessage("Failed to rename game: " + err.Error())
			errorDialog.Show()
			return
		}
		actionsMenu.UpdateGame(-1) // reload game info
		menu.UpdateItemList()
		actionsMenu.Show()
	})
	screen.root.AddButton("Cancel", actionsMenu.Show)
	pages.AddPage("actionsRename", screen.root, true, false)
	return screen
}

func (r *ActionsRenameScreen) Show() {
	pages.SwitchToPage("actionsRename")
}

func (r *ActionsRenameScreen) SetGameIndex(gameIndex int) {
	r.gameIndex = gameIndex
	game := manager.Get(gameIndex)
	r.root.Clear(false)
	r.root.AddTextView("[white]Old name:", "[purple]"+game.GetName(), 0, 1, true, false)
	r.newName = game.GetName()
	r.root.AddInputField(
		"New name:",
		r.newName,
		manager.MaxNameSize,
		manager.CheckIfAcceptName,
		func(t string) { r.newName = t },
	)
}
