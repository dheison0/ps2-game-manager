package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/rivo/tview"
)

func GameActions(app *tview.Application, pages *tview.Pages, game manager.GameConfig) *tview.List {
	actions := tview.NewList()
	actions.SetBorder(true)
	actions.SetTitle(fmt.Sprintf("Actions for '%s'", game.Name))
	actions.AddItem("Rename", "Rename game", 'r', func() {
		pages.AddAndSwitchToPage("Rename", renameAction(app, pages, game), true)
	})
	actions.AddItem("Delete", "Remove game from disk", 'd', func() {
		pages.AddAndSwitchToPage("Delete", deleteAction(app, pages, game), true)
	})
	actions.AddItem("Back", "Go to main menu", 'b', func() { pages.SwitchToPage("Menu") })
	return actions
}

func renameAction(app *tview.Application, pages *tview.Pages, game manager.GameConfig) *tview.Flex {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	text := tview.NewTextView()
	text.SetText(fmt.Sprintf("Name: %s\nImage: %s", game.Name, game.Image))
	flex.AddItem(text, 0, 1, false)
	form := tview.NewForm()
	newName := fmt.Sprintf("%s", game.Name)
	form.AddInputField("New name", newName, len(game.Name), func(t string, _ rune) bool { return len(t) <= len(game.Name) }, func(t string) { newName = t })
	form.AddButton("Save", func() {
		copy(game.Name[:], []byte(newName))
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Cancel", func() { pages.SwitchToPage("Menu") })
	flex.AddItem(form, 0, 1, true)
	return flex
}

func deleteAction(app *tview.Application, pages *tview.Pages, game manager.GameConfig) *tview.Modal {
	return nil
}
