package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/rivo/tview"
)

func GameActions(index int) tview.Primitive {
	game := manager.Get(index)
	actions := tview.NewList()
	actions.SetBorder(true)
	actions.SetTitle(fmt.Sprintf("Actions for '%s'", game.GetName()))
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
	game := manager.Get(index)
	newName := game.GetName()

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	text := tview.NewTextView()
	text.SetText(fmt.Sprintf("Name: %s\nImage: %s", game.GetName(), game.GetImage()))

	form := tview.NewForm()
	form.AddInputField(
		"New name:", newName, len(game.Config.Name),
		func(t string, _ rune) bool {
			return len(t) <= len(game.Config.Name)
		},
		func(t string) { newName = t },
	)
	form.AddButton("Save", func() {
		if err := manager.Rename(index, newName); err != nil {
			ErrorDialog("Failed to rename game to '" + newName + "'!")
		}
		Menu()
	})
	form.AddButton("Cancel", func() { pages.SwitchToPage("Menu") })

	flex.AddItem(text, 2, 1, false)
	flex.AddItem(form, 0, 1, true)
	return flex
}

func ActionDelete(index int) tview.Primitive {
	game := manager.Get(index)
	modal := tview.NewModal()
	modal.SetText(fmt.Sprintf("Do you really want to delete '%s' game?", game.GetName()))
	modal.AddButtons([]string{"Abort", "Confirm"})
	modal.SetDoneFunc(func(buttonID int, _ string) {
		if buttonID != 1 {
			pages.AddAndSwitchToPage("Game", GameActions(index), true)
		}
		err := manager.Delete(index)
		if err != nil {
			ErrorDialog(fmt.Sprintf("Failed to remove files or game from config:\n%v", err))
		}
		Menu()
	})
	return modal
}
