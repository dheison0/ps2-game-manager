package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/rivo/tview"
)

func Menu(app *tview.Application, pages *tview.Pages) *tview.List {
	games := manager.GetGames()
	list := tview.NewList()
	list.SetTitle("PS2 Game Manager")
	list.SetBorder(true)
	for i := range games {
		list.AddItem(fmt.Sprintf("%s", games[i].Name), "", rune(i), func() {
			game := games[list.GetCurrentItem()]
			gamePage := GameActions(app, pages, game)
			pages.AddAndSwitchToPage("Game", gamePage, true)
		})
	}
	list.AddItem("Quit", "Exit program", 'q', func() {
		app.Stop()
	})
	return list
}
