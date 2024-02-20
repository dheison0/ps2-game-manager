package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/rivo/tview"
)

func GameActions(index int) tview.Primitive {
	gameCfg := manager.GetGame(index).Config
	actions := tview.NewList()
	actions.SetBorder(true)
	actions.SetTitle(fmt.Sprintf("Actions for '%s'", gameCfg.Name))
	actions.AddItem("Rename", "Rename game", 'r', func() {
		pages.AddAndSwitchToPage("Rename", ActionRename(index), true)
	})
	actions.AddItem("Delete", "Remove game from disk", 'd', func() {
		pages.AddAndSwitchToPage("Delete", ActionDelete(index), true)
	})
	actions.AddItem("Back", "Go to main menu", 'b', func() { pages.SwitchToPage("Menu") })
	return actions
}

func ActionRename(index int) tview.Primitive {
	gameCfg := manager.GetGame(index).Config
	newName := string(gameCfg.Name[:])

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	text := tview.NewTextView()
	text.SetText(fmt.Sprintf("Name: %s\nImage: %s", gameCfg.Name, gameCfg.Image))

	form := tview.NewForm()
	form.AddInputField(
		"New name:", newName, len(gameCfg.Name),
		func(t string, _ rune) bool {
			return len(t) <= len(gameCfg.Name)
		},
		func(t string) { newName = t },
	)
	form.AddButton("Save", func() {
		gameCfg.Name = [len(gameCfg.Name)]byte{}
		copy(gameCfg.Name[:], []byte(newName))
		manager.UpdateGameConfig(index, gameCfg)
		RefreshPages()
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Cancel", func() { pages.SwitchToPage("Menu") })

	flex.AddItem(text, 2, 1, false)
	flex.AddItem(form, 0, 1, true)
	return flex
}

func ActionDelete(index int) tview.Primitive {
	gameCfg := manager.GetGame(index).Config
	modal := tview.NewModal()
	modal.SetText(fmt.Sprintf("Do you really want to delete '%s' game?", gameCfg.Name))
	modal.AddButtons([]string{"Abort", "Confirm"})
	modal.SetDoneFunc(func(buttonID int, _ string) {
		if buttonID != 1 {
			pages.AddAndSwitchToPage("Game", GameActions(index), true)
		}
		err := manager.RemoveGame(index)
		if err != nil {
			errorModal := tview.NewModal()
			errorModal.SetText(fmt.Sprintf("Failed to remove files or game from config:\n%v", err))
			errorModal.AddButtons([]string{"Done"})
			errorModal.SetDoneFunc(func(_ int, _ string) {
				pages.AddAndSwitchToPage("Game", GameActions(index), true)
			})
			pages.AddAndSwitchToPage("Error", errorModal, true)
		}
		RefreshPages()
		pages.SwitchToPage("Menu")

	})
	return modal
}
